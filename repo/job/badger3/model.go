package badger3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
)

type (
	modelJob struct {
		ID             jobq.ID         `json:"id"`
		Topic          jobq.Topic      `json:"topic"`
		Weight         jobq.Weight     `json:"pri"`
		Status         jobq.Status     `json:"status"`
		DateCreated    time.Time       `json:"dateCreated"`
		DateTerminated time.Time       `json:"dateTerminated"`
		DatesReserved  []time.Time     `json:"datesReserved"`
		RetryCount     uint            `json:"retryCount"`
		ResetCount     uint            `json:"resetCount"`
		Options        modelJobOptions `json:"options"`
		Logs           []string        `json:"logs"`
	}

	modelJobOptions struct {
		Name            string        `json:"name"`
		Timeout         time.Duration `json:"timeout"`
		DelayedAt       time.Time     `json:"delayedAt"`
		MaxRetries      uint          `json:"maxRetries"`
		MinBackOff      time.Duration `json:"minBackOff"`
		MaxBackOff      time.Duration `json:"maxBackOff"`
		LogStatusChange bool          `json:"LogStatusChange"`
	}
)

var (
	prefixKeyIdSequence  = []byte("j:id:seq:")
	prefixKeyJob         = []byte("j:j:")
	prefixKeyPayload     = []byte("j:p:")
	prefixKeyStatusIndex = []byte("i:j:s:")
)

func (m *modelJob) FromDomain(j *jobq.JobInfo) {
	m.ID = j.ID
	m.Topic = j.Topic
	m.Weight = j.Weight
	m.Status = j.Status
	m.DateCreated = j.DateCreated
	m.DateTerminated = j.DateTerminated
	m.DatesReserved = j.DatesReserved
	m.RetryCount = j.RetryCount
	m.ResetCount = j.ResetCount
	m.Logs = j.Logs
	m.Options.Name = j.Options.Name
	m.Options.Timeout = j.Options.Timeout
	m.Options.DelayedAt = j.Options.DelayedAt
	m.Options.MaxRetries = j.Options.MaxRetries
	m.Options.MinBackOff = j.Options.MinBackOff
	m.Options.MaxBackOff = j.Options.MaxBackOff
	m.Options.LogStatusChange = j.Options.LogStatusChange

}

func (m *modelJob) ToDomain() *jobq.JobInfo {
	return &jobq.JobInfo{
		ID:             m.ID,
		Topic:          m.Topic,
		Weight:         m.Weight,
		Status:         m.Status,
		DateCreated:    m.DateCreated,
		DateTerminated: m.DateTerminated,
		DatesReserved:  m.DatesReserved,
		RetryCount:     m.RetryCount,
		ResetCount:     m.ResetCount,
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

func (m *modelJob) keyID() []byte {
	return []byte(fmt.Sprintf("%020d", m.ID))
}

func (m *modelJob) keyJob() []byte {
	return append(prefixKeyJob, m.keyID()...)
}

func (m *modelJob) keyStatusIndex(status jobq.Status) []byte {
	key := new(bytes.Buffer)
	key.Write(prefixKeyStatusIndex)
	key.WriteString(status.String())
	key.WriteString(":")
	key.Write(m.keyID())

	return key.Bytes()
}

func (m *modelJob) keyPayload() []byte {
	return append(prefixKeyPayload, m.keyID()...)
}

func (m *modelJob) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *modelJob) Decode(data []byte) error {
	mjo := modelJob{}
	err := json.Unmarshal(data, &mjo)
	if err == nil {
		*m = mjo
	}

	return err
}
