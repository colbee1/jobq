package jqs

import (
	"context"

	"github.com/colbee1/jobq"
)

func (s *Service) Topics(ctx context.Context, offset int, limit int) ([]jobq.JobTopic, error) {
	return s.pqRepo.Topics(ctx, offset, limit)
}

// TopicStats returns some stats about topic.
func (s *Service) TopicStats(ctx context.Context, topic jobq.JobTopic) (jobq.TopicStats, error) {
	return s.pqRepo.TopicStats(ctx, topic)
}
