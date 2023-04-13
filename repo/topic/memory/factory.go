package memory

import (
	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/topic"
)

func New(options topic.RepositoryOptions) (*Adapter, error) {
	a := &Adapter{
		topics:  make(map[jobq.Topic]*JobQueue),
		delayed: newJobQueue(),
	}

	if options.StatsCollector {
		a.statsCollector = topic.StartStatsCollector(topic.CollectorDefaultQueueSize)
	}

	return a, nil
}

func (a *Adapter) Close() error {
	if a.statsCollector != nil {
		a.statsCollector.Close()
	}

	a.topics = make(map[jobq.Topic]*JobQueue)

	return nil
}
