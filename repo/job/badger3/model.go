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
		Priority       jobq.Priority   `json:"pri"`
		DateCreated    time.Time       `json:"dateCreated"`
		DateTerminated time.Time       `json:"dateTerminated"`
		Options        modelJobOptions `json:"options"`
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
	prefixKeyStatus      = []byte("j:s:")
	prefixKeyStatusIndex = []byte("i:j:s:")
	prefixKeyDateUpdated = []byte("j:du:")
	prefixKeyInfo        = []byte("j:i:")
	prefixKeyPayload     = []byte("j:p:")
	prefixKeyLogs        = []byte("j:l:")
)

func (m *modelJob) keyID() []byte {
	return []byte(fmt.Sprintf("%020d", m.ID))
}

func (m *modelJob) keyJob() []byte {
	return append(prefixKeyJob, m.keyID()...)
}

func (m *modelJob) keyStatus() []byte {
	return append(prefixKeyStatus, m.keyID()...)
}

func (m *modelJob) keyStatusIndex(status jobq.Status) []byte {
	key := new(bytes.Buffer)
	key.Write(prefixKeyStatusIndex)
	key.WriteString(status.String())
	key.WriteString(":")
	key.Write(m.keyID())

	return key.Bytes()
}

func (m *modelJob) keyUpdated() []byte {
	return append(prefixKeyDateUpdated, m.keyID()...)
}

func (m *modelJob) keyLogs() []byte {
	return append(prefixKeyLogs, m.keyID()...)
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
	if err != nil {
		return err
	}
	*m = mjo

	return nil
}
