package jqs

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/service"
)

// Reserve reserves up to <limit> Ready jobs
func (s *Service) Reserve(ctx context.Context, topic jobq.Topic, limit int) ([]service.IJob, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	if limit < 1 {
		return []service.IJob{}, nil
	}

	jobs := make([]service.IJob, 0, limit)

	reserve := func(max int) error {
		jids, err := s.pqRepo.PopTopic(ctx, topic, max)
		if err != nil {
			return err
		}
		if len(jids) == 0 {
			return nil
		}

		tx, err := s.jobRepo.NewTransaction()
		if err != nil {
			return err
		}
		defer tx.Close()
		err = tx.Update(context.Background(), jids,
			func(job *jobq.JobInfo) error {
				job.Status = jobq.JobStatusReserved
				job.DatesReserved = append(job.DatesReserved, time.Now())

				jobs = append(jobs, &Job{service: s, id: job.ID})

				return nil
			})
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	if err := reserve(limit); err != nil {
		return jobs, err
	}

	_, withTimeout := ctx.Deadline()
	if !withTimeout || len(jobs) == limit {
		return jobs, nil
	}

	fmt.Printf("reserve with timeout...\n")

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

loop:
	for {
		select {

		case <-ctx.Done():
			break loop

		case <-ticker.C:
			if err := reserve(limit - len(jobs)); err != nil {
				return jobs, err
			}
			if len(jobs) == limit {
				break loop
			}
		}
	}

	return jobs, nil
}
