package memory

import (
	"sync"
	"sync/atomic"

	"github.com/colbee1/jobq"
)

var _ jobq.IJobRepository = (*Adapter)(nil)

type Adapter struct {
	jobSequence atomic.Uint64 // same type as jobq.JobID
	jobs        map[jobq.JobID]*modelJob
	jobsLock    sync.RWMutex
}

func (a *Adapter) Close() error {
	a.jobs = make(map[jobq.JobID]*modelJob)

	return nil
}
