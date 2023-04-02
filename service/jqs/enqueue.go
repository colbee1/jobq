package jqs

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (s *Service) Enqueue(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	if jo.Timeout == 0 {
		jo.Timeout = jobq.DefaultJobOptions.Timeout
	}

	if jo.MaxRetries == 0 {
		jo.MaxRetries = jobq.DefaultJobOptions.MaxRetries
	}

	if jo.MinBackOff < time.Second {
		jo.MinBackOff = time.Second
	}

	if jo.MaxBackOff == 0 {
		jo.MaxBackOff = jobq.DefaultJobOptions.MaxBackOff
	}

	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return 0, err
	}
	defer tx.Close()

	jid, err := tx.Create(ctx, topic, pri, jo, payload)
	if err != nil {
		return 0, err
	}

	err = tx.Update(context.Background(), []jobq.ID{jobq.ID(jid)},
		func(job *jobq.JobInfo) error {
			status, err := s.pqRepo.Push(ctx, topic, pri, jid, jo.DelayedAt)
			if err != nil {
				return err
			}

			job.Status = status
			return nil
		})
	if err != nil {
		return 0, nil
	}

	return jid, tx.Commit()
}
