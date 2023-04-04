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

	if limit == 0 {
		return []service.IJob{}, nil
	}

	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	jids, err := s.pqRepo.PopTopic(ctx, topic, limit)
	if err != nil {
		return nil, err
	}
	fmt.Printf("____ poped jids: %+v\n", jids)

	if len(jids) == 0 {
		return []service.IJob{}, nil
	}

	jobs := make([]service.IJob, 0, limit)
	err = tx.Update(context.Background(), jids,
		func(job *jobq.JobInfo) error {
			job.Status = jobq.JobStatusReserved
			job.DatesReserved = append(job.DatesReserved, time.Now())

			jobs = append(jobs, &Job{service: s, id: job.ID})

			return nil
		})
	if err != nil {
		for _, job := range jobs {
			// TODO: find a better way to handle orphan jobs
			job.Retry(0) // 0 == use job backoff delay
		}

		return nil, err
	}

	return jobs, tx.Commit()
}
