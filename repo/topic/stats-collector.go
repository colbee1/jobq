package topic

import (
	"fmt"
	"sync"
	"time"

	"github.com/colbee1/assertor"
	"github.com/colbee1/jobq"
)

type Metric int

const (
	MetricPush Metric = iota
	MetricPop  Metric = iota
)

const CollectorDefaultQueueSize = 128

type StatMetric struct {
	Topic  jobq.Topic
	Metric Metric
	Value  int64
}

type topicStats struct {
	pushDateFirst  time.Time
	pushDateLast   time.Time
	pushTotalCount int64
	popDateFirst   time.Time
	popDateLast    time.Time
	popTotalCount  int64
}

type StatsCollector struct {
	Collector   chan StatMetric
	wg          *sync.WaitGroup
	dateCreated time.Time
	byTopics    map[jobq.Topic]*topicStats
}

func StartStatsCollector(qsize int) *StatsCollector {
	sc := &StatsCollector{
		Collector:   make(chan StatMetric, qsize),
		wg:          new(sync.WaitGroup),
		dateCreated: time.Now(),
		byTopics:    make(map[jobq.Topic]*topicStats),
	}

	sc.wg.Add(1)
	go sc.collect()

	return sc
}

func (sc *StatsCollector) Topics(offset, limit int64) ([]jobq.Topic, error) {
	v := assertor.New()
	v.Assert(limit > 0, "invalid limit")
	if err := v.Validate(); err != nil {
		return nil, err
	}

	ret := []jobq.Topic{}
	for t := range sc.byTopics {
		if offset > 0 {
			offset--
			continue
		}

		ret = append(ret, t)

		limit--
		if limit == 0 {
			break
		}

	}

	return ret, nil
}

func (sc *StatsCollector) TopicStats(topic jobq.Topic) (Stats, error) {
	stats := sc.byTopics[topic]
	if stats == nil {
		return Stats{}, ErrTopicNotFound
	}

	ret := Stats{
		PushDateFirst:  stats.pushDateFirst,
		PushDateLast:   stats.pushDateLast,
		PushTotalCount: stats.pushTotalCount,
		PopDateFirst:   stats.popDateFirst,
		PopDateLast:    stats.popDateLast,
		PopTotalCount:  stats.popTotalCount,
	}

	return ret, nil
}

func (sc *StatsCollector) Close() error {
	close(sc.Collector)
	sc.wg.Wait()

	return nil
}

func (sc *StatsCollector) collect() {
	defer sc.wg.Done()

	for metric := range sc.Collector {

		now := time.Now()
		t := metric.Topic
		stats := sc.byTopics[t]
		if stats == nil {
			stats = new(topicStats)
			sc.byTopics[t] = stats
		}

		switch metric.Metric {

		case MetricPush:
			stats.pushTotalCount += metric.Value
			if stats.pushDateFirst.IsZero() {
				stats.pushDateFirst = now
			}
			stats.pushDateLast = now
			// fmt.Printf("metric PUSH: t=%s, v=%d, total=%d, qlen=%d\n", t, metric.Value, stats.pushTotalCount, stats.pushTotalCount-stats.popTotalCount)

		case MetricPop:
			stats.popTotalCount += metric.Value
			if stats.popDateFirst.IsZero() {
				stats.popDateFirst = now
			}
			stats.popDateLast = now
			// fmt.Printf("metric POP: t=%s, v=%d, total=%d, qlen=%d\n", t, metric.Value, stats.popTotalCount, stats.pushTotalCount-stats.popTotalCount)

		default:
			fmt.Printf("unknown metric: %v\n", metric.Metric)
		}
	}
}
