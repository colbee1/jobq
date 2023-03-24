package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (a *Adapter) Topics(ctx context.Context, offset int, limit int) ([]jobq.JobTopic, error) {
	ret := make([]jobq.JobTopic, 0, len(a.pqByTopic))

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

func (a *Adapter) TopicStats(ctx context.Context, topic jobq.JobTopic) (jobq.TopicStats, error) {
	if pq, found := a.pqByTopic[topic]; !found {
		return jobq.TopicStats{}, jobq.ErrTopicNotFound
	} else {
		return pq.Stats(), nil
	}
}
