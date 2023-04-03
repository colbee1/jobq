package jqs

import (
	"sync"

	"github.com/colbee1/jobq/repo"
	"github.com/colbee1/jobq/service"
)

var _ service.IJobQueueService = (*Service)(nil)

const DefaultTopic = "$dflt"

type Service struct {
	jobRepo   repo.IJobRepository
	pqRepo    repo.IJobPriorityQueueRepository
	wg        sync.WaitGroup
	exitSched chan struct{}
}
