package memory

import (
	"sync"
	"sync/atomic"

	"github.com/colbee1/jobq"
)

type Adapter struct {
	jobSequence atomic.Uint64 // same type as jobq.JobID
	jobs        map[jobq.JobID]*jobq.Job
	jobsLock    sync.RWMutex
}
