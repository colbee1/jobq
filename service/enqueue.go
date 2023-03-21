package service

import (
	"context"

	"github.com/colbee1/jobq"
)

func (s *Service) Enqueue(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jo jobq.JobOptions) (jobq.JobID, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	jid, err := s.jobRepo.Create(ctx, topic, pri, jo)
	if err != nil {
		return 0, err
	}

	if err := s.pqRepo.Push(ctx, topic, pri, jid, jo.DelayedAt); err != nil {
		return 0, err
	}

	return jid, nil
}
