package memory

import (
	"context"
	"time"

	"github.com/colbee1/assertor"

	"github.com/colbee1/jobq"
)

func (a *Adapter) Push(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jid jobq.JobID, delayedAt time.Time) (jobq.JobStatus, error) {
	v := assertor.New()
	v.Assert(ctx != nil, "context is missing")
	v.Assert(topic != "", "topic is missing")
	if err := v.Validate(); err != nil {
		return 0, err
	}

	jitem := &JobItem{
		Topic:    topic,
		Priority: pri,
		JobID:    jid,
	}

	if time.Until(delayedAt) > 1 {
		jitem.heapPriority = delayedAt.Unix()
		a.pqDelayed.Push(jitem)

		return jobq.JobStatusDelayed, nil
	}

	pq, found := a.pqByTopic[topic]
	if !found {
		a.pqByTopic[topic] = newJobQueue()
		pq = a.pqByTopic[topic]
	}

	jitem.heapPriority = int64(jitem.Priority)
	pq.Push(jitem)

	return jobq.JobStatusReady, nil
}
