package jobq

import (
	"time"
)

type (
	Job struct {
		ID          JobID
		Status      JobStatus
		Priority    JobPriority
		Topic       JobTopic
		DateCreated time.Time
		DateUpdated time.Time
		RetryCount  int
		Date        time.Time
		Message     string
		Options     JobOptions
	}

	JobOptions struct {
		Name            string
		DelayedAt       time.Time
		RetryBackoff    time.Duration
		RetryMaxBackoff time.Duration
		RetryLimit      int
		Payload         JobPayload
	}

	JobID       uint64
	JobStatus   int16
	JobPriority int16
	JobTopic    string
	JobPayload  []byte
)

const (
	JobStateBuried   = JobStatus(-2)
	JobStateFail     = JobStatus(-1)
	JobStateReady    = JobStatus(0)
	JobStateReserved = JobStatus(1)
	JobStateDone     = JobStatus(2)
)
