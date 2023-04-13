package jqs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/job"
	"github.com/colbee1/jobq/repo/topic"
	"github.com/colbee1/jobq/service"
)

var _ service.IJobService = (*Service)(nil)

const DefaultTopic = "default"

type Service struct {
	jobRepo       job.IJobRepository
	topicRepo     topic.ITopicRepository
	wg            sync.WaitGroup
	stopScheduler chan struct{}
}

func (s *Service) Enqueue(ctx context.Context, topic jobq.Topic, pri jobq.Weight, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
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

	err = tx.Update(context.Background(), []jobq.ID{jid},
		func(job *jobq.JobInfo) error {
			status, err := s.topicRepo.Push(ctx, topic, pri, jid, jo.DelayedAt)
			if err != nil {
				return err
			}
			job.Status = status

			return nil
		})
	if err != nil {
		return 0, err
	}

	return jid, tx.Commit()
}

func (s *Service) Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error) {
	return s.topicRepo.Topics(ctx, offset, limit)
}

// TopicStats returns some stats about topic.
func (s *Service) TopicStats(ctx context.Context, topic jobq.Topic) (topic.Stats, error) {
	return s.topicRepo.TopicStats(ctx, topic)
}

// Reset resets a job and requeue it.
// RetryCount and DateTerminated are reset.
// ResetCount is incremented.
// Job is immediately pushed in it's topic queue.
func (s *Service) Reset(ctx context.Context, jids []jobq.ID) error {
	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	err = tx.Update(ctx, jids, func(job *jobq.JobInfo) error {
		status, err := s.topicRepo.Push(ctx, job.Topic, job.Weight, job.ID, time.Time{})
		if err != nil {
			return err
		}

		job.RetryCount = 0
		job.ResetCount++
		job.Status = status
		job.DateTerminated = time.Time{}

		return nil
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

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
		jids, err := s.topicRepo.PopTopic(ctx, topic, max)
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
				dateReserved := time.Now()

				job.Status = jobq.JobStatusReserved
				job.DatesReserved = append(job.DatesReserved, dateReserved)

				jobs = append(jobs, newJob(s, job.ID, dateReserved))

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

	ticker := time.NewTicker(250 * time.Millisecond) // TODO: make it configurable
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

func (s *Service) FindByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error) {
	tx, err := s.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	return tx.FindByStatus(ctx, status, offset, limit)
}
