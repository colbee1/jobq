package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (a *Adapter) PopTopic(_ context.Context, topic jobq.Topic, limit int) ([]jobq.ID, error) {
	if limit < 1 {
		return []jobq.ID{}, nil
	}

	if topic == "" {
		return nil, jobq.ErrTopicIsMissing
	}

	pq, found := a.pqByTopic[topic]
	if !found {
		a.pqByTopic[topic] = newJobQueue()
		pq = a.pqByTopic[topic]
	}

	jobs, err := pq.Pop(limit)
	if err != nil {
		return nil, err
	}
	jids := make([]jobq.ID, 0, len(jobs))
	for _, job := range jobs {
		jids = append(jids, job.JobID)
	}

	return jids, nil
}

func (a *Adapter) PopDelayed(_ context.Context, limit int) ([]jobq.ID, error) {
	if limit < 1 {
		return []jobq.ID{}, nil
	}

	now := time.Now().Unix()
	jids := make([]jobq.ID, 0, limit)
	for {
		jitem := a.pqDelayed.Peek()
		if jitem == nil || jitem.heapPriority > now {
			break
		}

		jitems, _ := a.pqDelayed.Pop(1)
		if len(jitems) == 1 {
			jids = append(jids, jitems[0].JobID)
		}

		if len(jids) == limit {
			break
		}
	}

	return jids, nil
}
