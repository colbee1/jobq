package jqs

import (
	"context"
	"fmt"

	"github.com/colbee1/jobq"
)

func (s *Service) Available(ctx context.Context, topic jobq.Topic) (int, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	return 0, fmt.Errorf("not yet available")
}

func (s *Service) Delayed(ctx context.Context) (int, error) {
	return 0, fmt.Errorf("not yet available")
}
