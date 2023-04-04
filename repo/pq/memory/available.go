package memory

import (
	"context"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

func (a *Adapter) AvailableTopic(ctx context.Context, topic jobq.Topic) (int, error) {
	pq, found := a.pqByTopic[topic]
	if !found {
		return -1, repo.ErrTopicNotFound
	}

	return pq.Len(), nil
}

func (a *Adapter) AvailableDelayed(ctx context.Context) (int, error) {
	return a.pqDelayed.Len(), nil
}
