package memory

import (
	"sync"
	"sync/atomic"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

var _ repo.IJobRepository = (*Adapter)(nil)

type Adapter struct {
	jobSequence atomic.Uint64 // same type as jobq.JobID
	jobs        map[jobq.ID]*modelJob
	jobsLock    sync.RWMutex
}
