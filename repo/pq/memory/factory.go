package memory

import (
	"github.com/colbee1/jobq"
)

func New() (*Adapter, error) {
	a := &Adapter{
		pqByTopic: make(map[jobq.Topic]*JobQueue),
		pqDelayed: newJobQueue(),
	}

	return a, nil
}

func (a *Adapter) Close() error {
	a.pqByTopic = make(map[jobq.Topic]*JobQueue)

	return nil
}
