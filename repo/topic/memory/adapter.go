package memory

import (
	"context"
	"time"

	"github.com/colbee1/assertor"
	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/topic"
)

type Adapter struct {
	topics         map[jobq.Topic]*JobQueue
	delayed        *JobQueue
	statsCollector *topic.StatsCollector
}

func (a *Adapter) createTopic(ctx context.Context, topic jobq.Topic) *JobQueue {
	pq, found := a.topics[topic]
	if !found {
		pq = newJobQueue()
		a.topics[topic] = pq
	}

	return pq
}

func (a *Adapter) Push(ctx context.Context, tp jobq.Topic, pri jobq.Weight, jid jobq.ID, delayedAt time.Time) (jobq.Status, error) {
	v := assertor.New()
	v.Assert(ctx != nil, "context is missing")
	v.Assert(tp != "", "topic is missing")
	if err := v.Validate(); err != nil {
		return 0, err
	}

	var pq *JobQueue
	var status jobq.Status
	jitem := &JobItem{
		Topic:    tp,
		Priority: pri,
		JobID:    jid,
	}

	if time.Until(delayedAt) > time.Second {
		jitem.heapPriority = delayedAt.Unix()
		status = jobq.JobStatusDelayed
		pq = a.delayed
	} else {
		jitem.heapPriority = int64(jitem.Priority)
		status = jobq.JobStatusReady
		pq = a.createTopic(ctx, tp)
	}

	pq.Push(jitem)
	if a.statsCollector != nil {
		t := tp
		if status == jobq.JobStatusDelayed {
			t = "delayed"
		}
		a.statsCollector.Collector <- topic.StatMetric{Topic: t, Metric: topic.MetricPush, Value: 1}
	}

	return status, nil
}

func (a *Adapter) PopTopic(ctx context.Context, tp jobq.Topic, limit int) ([]jobq.ID, error) {
	if limit < 1 {
		return []jobq.ID{}, nil
	}

	if tp == "" {
		return nil, jobq.ErrTopicIsMissing
	}

	pq := a.createTopic(ctx, tp)
	jids := make([]jobq.ID, 0, limit)
	for limit > 0 {
		job := pq.Pop()
		if job == nil {
			break
		}
		limit--

		jids = append(jids, job.JobID)
		if a.statsCollector != nil {
			a.statsCollector.Collector <- topic.StatMetric{Topic: tp, Metric: topic.MetricPop, Value: 1}
		}
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
		jitem := a.delayed.Peek()
		if jitem == nil || jitem.heapPriority > now {
			break
		}

		job := a.delayed.Pop()
		if job != nil {
			jids = append(jids, job.JobID)

			if a.statsCollector != nil {
				// TODO: Fix static topic
				a.statsCollector.Collector <- topic.StatMetric{Topic: "delayed", Metric: topic.MetricPop, Value: 1}
			}
		}

		if len(jids) == limit {
			break
		}
	}

	return jids, nil
}

func (a *Adapter) Durable() bool {
	return false
}

func (a *Adapter) Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error) {
	if a.statsCollector == nil {
		return nil, topic.ErrStatsCollectorDisabled
	}

	return a.statsCollector.Topics(int64(offset), int64(limit))
}

func (a *Adapter) TopicStats(ctx context.Context, tp jobq.Topic) (topic.Stats, error) {
	if a.statsCollector == nil {
		return topic.Stats{}, topic.ErrStatsCollectorDisabled
	}

	return a.statsCollector.TopicStats(tp)
}
