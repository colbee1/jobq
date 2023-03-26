package badger3

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
	"github.com/dgraph-io/badger/v3"
)

func (t *Transaction) Create(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
	a := t.a

	if topic == "" {
		return 0, jobq.ErrTopicIsMissing
	}

	newID := jobq.ID(0)
	for newID == 0 {
		if id, err := a.jobIdSeq.Next(); err != nil {
			return 0, err
		} else {
			newID = jobq.ID(id)
		}
	}

	model := &modelJob{
		ID:          newID,
		Topic:       topic,
		Priority:    pri,
		DateCreated: time.Now(),
		Options: modelJobOptions{
			Name:            jo.Name,
			Timeout:         jo.Timeout,
			DelayedAt:       jo.DelayedAt,
			MaxRetries:      jo.MaxRetries,
			MinBackOff:      jo.MinBackOff,
			MaxBackOff:      jo.MaxBackOff,
			LogStatusChange: jo.LogStatusChange,
		},
	}

	jobData, err := model.Encode()
	if err != nil {
		return 0, err
	}

	txn := t.tx
	if err := txn.Set(model.keyJob(), jobData); err != nil {
		return 0, err
	}
	if err := txn.Set(model.keyStatus(), []byte(fmt.Sprint(jobq.JobStatusCreated))); err != nil {
		return 0, err
	}
	if err := txn.Set(model.keyPayload(), payload); err != nil {
		return 0, err
	}

	if model.Options.LogStatusChange {
		if err := t.Log(ctx, newID, "job created"); err != nil {
			return 0, err
		}
	}

	t.needCommit = true

	return newID, nil
}

func (t *Transaction) SetStatus(ctx context.Context, jids []jobq.ID, status jobq.Status) error {
	m := &modelJob{}

	for _, jid := range jids {
		m.ID = jid
		// get current status to removefrom index
		//
		itm, err := t.tx.Get(m.keyStatus())
		if err == nil {
			data, err := itm.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := t.tx.Delete(m.keyStatusIndex(jobq.NewStatusFromString(string(data)))); err != nil {
				return err
			}
		} else if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		if err := t.tx.Set(m.keyStatus(), []byte(status.String())); err != nil {
			return err
		}
		if err := t.tx.Set(m.keyStatusIndex(status), []byte{}); err != nil {
			return err
		}

		t.needCommit = true
	}

	return nil
}

func (t *Transaction) Log(ctx context.Context, jid jobq.ID, log string) error {
	if log == "" {
		return nil
	}

	logs, err := t.Logs(ctx, jid)
	if err != nil {
		return err
	}

	logs = append(logs, fmt.Sprintf("%s: %s", time.Now().Format(time.RFC3339Nano), log))
	data, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	m := &modelJob{ID: jid}
	err = t.tx.Set(m.keyLogs(), data)
	if err != nil {
		return err
	}

	t.needCommit = true

	return nil
}

func (t *Transaction) SetOptions(ctx context.Context, jid jobq.ID, jo *jobq.JobOptions) error {
	return fmt.Errorf("not yet implemented") //TODO:
}
