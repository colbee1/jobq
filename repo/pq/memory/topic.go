package memory

import (
	"context"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

func (a *Adapter) Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error) {
	ret := make([]jobq.Topic, 0, len(a.pqByTopic))

	if limit < 1 {
		limit = 1000
	}

	for t := range a.pqByTopic {
		if offset > 0 {
			offset--
			continue
		}

		ret = append(ret, t)
		if len(ret) == limit {
			break
		}
	}

	return ret, nil
}

func (a *Adapter) TopicStats(ctx context.Context, topic jobq.Topic) (repo.TopicStats, error) {
	if pq, found := a.pqByTopic[topic]; !found {
		return repo.TopicStats{}, jobq.ErrTopicNotFound
	} else {
		return pq.Stats(), nil
	}
}
