package service

import (
	"context"

	"github.com/colbee1/jobq"
)

// Reserve reserves up to <limit> ready jobs
func (s *Service) Reserve(ctx context.Context, topic jobq.JobTopic, limit int) ([]*Job, error) {
	if topic == "" {
		topic = DefaultTopic
	}

	// Pop ids from priority queue
	//
	jids, err := s.pqRepo.Pop(ctx, topic, limit)
	if err != nil {
		return nil, err
	}
	if len(jids) == 0 {
		return []*Job{}, nil
	}

	// Change jobs state
	//
	if err := s.jobRepo.SetStatus(ctx, jids, jobq.JobStateReserved); err != nil {
		// TODO: how to handle poped job ? retry later ?
		return nil, err
	}

	jobs := make([]*Job, len(jids))
	for i, jid := range jids {
		jobs[i] = &Job{
			service: s,
			jid:     jid,
			topic:   topic,
		}
	}

	return jobs, nil
}
