package jqs

import (
	"context"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

func (s *Service) Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error) {
	return s.pqRepo.Topics(ctx, offset, limit)
}

// TopicStats returns some stats about topic.
func (s *Service) TopicStats(ctx context.Context, topic jobq.Topic) (repo.TopicStats, error) {
	return s.pqRepo.TopicStats(ctx, topic)
}
