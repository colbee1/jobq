package jqs

import (
	"context"

	"github.com/colbee1/jobq"
)

// Reserve reserves up to <limit> ready jobs
func (s *Service) Reserve(ctx context.Context, topic jobq.JobTopic, limit int) ([]jobq.IJob, error) {
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
		return []jobq.IJob{}, nil
	}

	// Change jobs state
	//
	if err := tx.SetStatus(ctx, jids, jobq.JobStatusReserved); err != nil {
		// TODO: handler poped job by re-push with delay ?
		return nil, err
	}

	jobs := make([]jobq.IJob, len(jids))
	for i, jid := range jids {
		jobs[i] = &Job{
			service: s,
			id:      jid,
			topic:   topic,
		}
	}

	return jobs, tx.Commit()
}
