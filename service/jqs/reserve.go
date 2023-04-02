package jqs

import (
	"context"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/service"
)

// Reserve reserves up to <limit> ready jobs
func (s *Service) Reserve(ctx context.Context, topic jobq.Topic, limit int) ([]service.IJobService, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	// Pop ids from priority queue
	//
	jids, err := s.pqRepo.Pop(ctx, topic, limit)
	if err != nil {
		return nil, err
	}
	if len(jids) == 0 {
		return []service.IJobService{}, nil
	}

	jobs := make([]service.IJobService, 0, len(jids))
	err = tx.Update(context.Background(), jids,
		func(job *jobq.JobInfo) error {
			job.Status = jobq.JobStatusReserved
			jobs = append(jobs, &Job{service: s, id: job.ID})
			return nil
		})
	if err != nil {
		for _, job := range jobs {
			// TODO: find a better way to handle orphan jobs
			job.Retry(0)
		}

		return nil, err
	}

	return jobs, tx.Commit()
}
