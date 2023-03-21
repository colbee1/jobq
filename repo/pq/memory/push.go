package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/assertor"

	"github.com/colbee1/jobq"
)

func (a *Adapter) Push(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jid jobq.JobID, delayedAt time.Time) error {
	v := assertor.New()
	v.Assert(ctx != nil, "context is missing")
	v.Assert(topic != "", "topic is missing")
	if err := v.Validate(); err != nil {
		return err
	}

	jitem := &JobItem{
		Topic:    topic,
		Priority: pri,
		JobID:    jid,
	}

	if !delayedAt.IsZero() {
		delay := time.Until(delayedAt)
		if delay > time.Second {

			ts := delayedAt.Unix()
			jitem.heapPriority = ts
			a.pqDelayed.Push(jitem)
			fmt.Printf("jid=%d is delayed at ts=%v\n", jitem.JobID, jitem.heapPriority)

			return nil
		}
	}

	pq, found := a.pqByTopic[topic]
	if !found {
		a.pqByTopic[topic] = newJobQueue()
		pq = a.pqByTopic[topic]
	}

	jitem.heapPriority = int64(jitem.Priority)
	pq.Push(jitem)

	return nil
}
