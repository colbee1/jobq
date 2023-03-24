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

func (a *Adapter) Len(ctx context.Context, topic jobq.JobTopic) (int, error) {
	pq, found := a.pqByTopic[topic]
	if !found {
		return 0, nil
	}

	return pq.Len(), nil
}
