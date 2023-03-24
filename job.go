package jobq

import (
	"time"
)

type IJob interface {
	ID() JobID
	Topic() JobTopic
	Priority() (JobPriority, error)
	Status() (JobStatus, error)
	// Payload() ([]byte, error)

	Log(string) error
	Logs() ([]string, error)

	// Options() (*JobOptions, error)

	Done(log string) error
	Fail(log string) error
	Cancel(log string) error
}

type JobOptions struct {
	Name            string
	Timeout         time.Duration
	DelayedAt       time.Time
	MaxRetries      uint
	MinBackOff      time.Duration
	MaxBackOff      time.Duration
	LogStatusChange bool
}

var DefaultJobOptions = JobOptions{
	Timeout:    2 * time.Hour,
	MaxRetries: 15,
	MinBackOff: 1 * time.Minute,
	MaxBackOff: 60 * time.Minute,
}

type JobInfo struct {
	ID             JobID
	Topic          JobTopic
	Priority       JobPriority
	Status         JobStatus
	DateCreated    time.Time
	DateTerminated time.Time
	DateReserved   []time.Time
	Retries        uint // == len(DateReserved) - 1
	Options        JobOptions
	Logs           []string
}

type (
	JobID       uint64
	JobStatus   int16
	JobPriority int16
	JobTopic    string
	JobPayload  []byte
)

const (
	JobStatusCanceled = JobStatus(-2)
	JobStatusDelayed  = JobStatus(-1)
	JobStatusCreated  = JobStatus(0)
	JobStatusReady    = JobStatus(1)
	JobStatusReserved = JobStatus(2)
	JobStatusDone     = JobStatus(3)
)

func (s JobStatus) String() string {
	switch s {
	case -2:
		return "Canceled"
	case -1:
		return "Delayed"
	case 0:
		return "Created"
	case 1:
		return "Ready"
	case 2:
		return "Reserved"
	case 3:
		return "Done"
	}

	return "Undefined!"
}
