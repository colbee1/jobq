package jobq

import (
	"math"
	"strconv"
	"strings"
	"time"
)

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
	MaxRetries: 50, // Max attempts before auto Cancel()
	MinBackOff: 30 * time.Second,
	MaxBackOff: 60 * time.Minute,
}

type JobInfo struct {
	ID             ID
	Topic          Topic
	Weight         Weight
	Status         Status
	DateCreated    time.Time
	DateTerminated time.Time
	DatesReserved  []time.Time
	RetryCount     uint
	ResetCount     uint
	Options        JobOptions
	Logs           []string
}

type (
	ID      uint64
	Status  int16
	Weight  int16
	Topic   string
	Payload []byte
)

const WeightMin = math.MinInt16
const WeightMax = math.MaxInt16

const (
	JobStatusUndefined = Status(-666)
	JobStatusCanceled  = Status(-2)
	JobStatusDelayed   = Status(-1)
	JobStatusCreated   = Status(0)
	JobStatusReady     = Status(1)
	JobStatusReserved  = Status(2)
	JobStatusDone      = Status(3)
)

func (s Status) String() string {
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

func NewTopicFromString(s string) Topic {
	return Topic(s)
}

func NewJobIDFromString(s string) (ID, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return ID(id), nil
	}

	return 0, err
}

func NewWeightFromString(s string) (Weight, error) {
	w, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return Weight(w), nil
	}

	return 0, err
}

func NewStatusFromString(s string) Status {
	switch strings.ToLower(s) {
	case "canceled":
		return JobStatusCanceled
	case "delayed":
		return JobStatusDelayed
	case "created":
		return JobStatusCreated
	case "ready":
		return JobStatusReady
	case "reserved":
		return JobStatusReserved
	case "done":
		return JobStatusDone
	}

	return JobStatusUndefined
}
