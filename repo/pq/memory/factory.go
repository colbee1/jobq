package memory

import (
	"github.com/colbee1/jobq"
)

func New() (*Adapter, error) {
	a := &Adapter{
		pqByTopic: make(map[jobq.JobTopic]*JobQueue),
		pqDelayed: newJobQueue(),
		exit:      make(chan struct{}, 1),
	}

	a.wg.Add(1)
	go a.delayedScheduler(a.exit)

	return a, nil
}

func (a *Adapter) Close() error {
	// Stop the scheduler for delayed jobs
	a.exit <- struct{}{}
	a.wg.Wait()

	// Free memory
	a.pqByTopic = make(map[jobq.JobTopic]*JobQueue)

	return nil
}
