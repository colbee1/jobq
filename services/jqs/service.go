package jqs

import (
	"github.com/colbee1/jobq"
)

var _ jobq.IJobQueueService = (*Service)(nil)

const DefaultTopic = "$dflt"

type Service struct {
	jobRepo jobq.IJobRepository
	pqRepo  jobq.IJobPriorityQueueRepository
}
