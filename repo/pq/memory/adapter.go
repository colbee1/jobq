package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

type Adapter struct {
	pqByTopic map[jobq.Topic]*JobQueue
	pqDelayed *JobQueue
}

func (a *Adapter) Len(ctx context.Context, topic jobq.Topic) (int, error) {
	pq, found := a.pqByTopic[topic]
	if !found {
		return 0, nil
	}

	return pq.Len(), nil
}
