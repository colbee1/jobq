package jqs

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (s *Service) Enqueue(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jo jobq.JobOptions, payload jobq.JobPayload) (jobq.JobID, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	if jo.Timeout == 0 {
		jo.Timeout = jobq.DefaultJobOptions.Timeout
	}

	if jo.MaxRetries == 0 {
		jo.MaxRetries = jobq.DefaultJobOptions.MaxRetries
	}

	if jo.MinBackOff < 5*time.Second {
		jo.MinBackOff = jobq.DefaultJobOptions.MinBackOff
	}

	if jo.MaxBackOff == 0 {
		jo.MaxBackOff = jobq.DefaultJobOptions.MaxBackOff
	}

	txJob, err := s.jobRepo.NewTransaction()
	if err != nil {
		return 0, err
	}
	defer txJob.Close()

	jid, err := txJob.Create(ctx, topic, pri, jo, payload)
	if err != nil {
		return 0, err
	}

	status, err := s.pqRepo.Push(ctx, topic, pri, jid, jo.DelayedAt)
	if err != nil {
		return 0, err
	}

	if err := txJob.SetStatus(ctx, []jobq.JobID{jid}, status); err != nil {
		return 0, err
	}

	return jid, txJob.Commit()
}
