package memory

import (
	"context"
	"sync"

	"github.com/colbee1/jobq"
)

type Adapter struct {
	pqByTopic map[jobq.JobTopic]*JobQueue
	pqDelayed *JobQueue
	exit      chan struct{}
	wg        sync.WaitGroup
}

func (a *Adapter) Pop(ctx context.Context, topic jobq.JobTopic, limit int) ([]jobq.JobID, error) {
	if topic == "" {
		return nil, jobq.ErrTopicIsInvalid
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

func (a *Adapter) Len(ctx context.Context, topic jobq.JobTopic) (int, error) {
	pq, found := a.pqByTopic[topic]
	if !found {
		return 0, nil
	}

	return pq.Len(), nil
}
