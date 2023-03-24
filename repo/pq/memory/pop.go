package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (a *Adapter) Pop(ctx context.Context, topic jobq.JobTopic, limit int) ([]jobq.JobID, error) {
	if topic == "" {
		return nil, jobq.ErrTopicIsMissing
	}

	pq, found := a.pqByTopic[topic]
	if !found {
		return nil, jobq.ErrTopicNotFound
	}

	if limit < 1 {
		limit = 1
	}

	jobs, err := pq.Pop(uint(limit))
	if err != nil {
		return nil, err
	}

	jids := make([]jobq.JobID, len(jobs))
	for i, job := range jobs {
		jids[i] = job.JobID
	}

	return jids, nil
}
