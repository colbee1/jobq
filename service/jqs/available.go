package jqs

import (
	"context"

	"github.com/colbee1/jobq"
)

func (s *Service) Available(ctx context.Context, topic jobq.Topic) (int, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	return s.pqRepo.AvailableTopic(ctx, topic)
}

func (s *Service) Delayed(ctx context.Context) (int, error) {
	return s.pqRepo.AvailableDelayed(ctx)
}
