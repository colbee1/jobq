package memory

import (
	"time"

	"github.com/colbee1/jobq"
)

type (
	modelJob struct {
		ID             jobq.ID
		Topic          jobq.Topic
		Priority       jobq.Priority
		Status         jobq.Status
		DateCreated    time.Time
		DateTerminated time.Time
		DatesReserved  []time.Time
		RetryCount     uint
		Options        modelJobOptions
		Logs           []string
		Payload        jobq.Payload
	}

	modelJobOptions struct {
		Name            string
		Timeout         time.Duration
		DelayedAt       time.Time
		MaxRetries      uint
		MinBackOff      time.Duration
		MaxBackOff      time.Duration
		LogStatusChange bool
	}
)

func (m *modelJob) ToDomain() *jobq.JobInfo {
	return &jobq.JobInfo{
		ID:             m.ID,
		Topic:          m.Topic,
		Priority:       m.Priority,
		Status:         m.Status,
		DateCreated:    m.DateCreated,
		DateTerminated: m.DateTerminated,
		DatesReserved:  m.DatesReserved,
		RetryCount:     m.RetryCount,
		Options: jobq.JobOptions{
			Name:            m.Options.Name,
			Timeout:         m.Options.Timeout,
			DelayedAt:       m.Options.DelayedAt,
			MaxRetries:      m.Options.MaxRetries,
			MinBackOff:      m.Options.MinBackOff,
			MaxBackOff:      m.Options.MaxBackOff,
			LogStatusChange: m.Options.LogStatusChange,
		},
		Logs: m.Logs,
	}
}

// type modelJob struct {
// 	id       jobq.ID
// 	topic    jobq.Topic
// 	priority jobq.Priority
// 	status   jobq.Status
// 	payload  jobq.Payload
// 	options  jobq.JobOptions
// 	info     jobq.JobInfo
// 	logs     []string
// }

// func (mj *modelJob) ToDomain() *jobq.JobInfo {
// 	return &jobq.JobInfo{
// 		ID:             mj.id,
// 		Topic:          mj.topic,
// 		Priority:       mj.priority,
// 		Status:         mj.status,
// 		DateCreated:    mj.info.DateCreated,
// 		DateTerminated: mj.info.DateTerminated,
// 		DatesReserved:  mj.info.DatesReserved,
// 		RetryCount:     mj.info.RetryCount,
// 		Options: jobq.JobOptions{
// 			Name:            mj.options.Name,
// 			Timeout:         mj.options.Timeout,
// 			DelayedAt:       mj.options.DelayedAt,
// 			MaxRetries:      mj.options.MaxRetries,
// 			MinBackOff:      mj.options.MinBackOff,
// 			MaxBackOff:      mj.options.MaxBackOff,
// 			LogStatusChange: mj.options.LogStatusChange,
// 		},
// 		Logs: mj.logs,
// 	}
// }
