package badger3

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

func (t *Transaction) saveJob(mj *modelJob) error {
	data, err := mj.Encode()
	if err != nil {
		return err
	}

	if err := t.tx.Set(mj.keyJob(), data); err != nil {
		return err
	}
	t.needCommit = true

	return nil
}

func (t *Transaction) Create(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
	if topic == "" {
		return 0, jobq.ErrTopicIsMissing
	}

	newID := jobq.ID(0)
	for newID == 0 {
		if id, err := t.a.jobIdSeq.Next(); err != nil {
			return 0, err
		} else {
			newID = jobq.ID(id)
		}
	}

	mj := &modelJob{
		ID:          newID,
		Topic:       topic,
		Priority:    pri,
		Status:      jobq.JobStatusCreated,
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

	if mj.Options.LogStatusChange {
		mj.Logs = append(mj.Logs, fmt.Sprintf("%s: Job status changed to: %s",
			time.Now().Format(time.RFC3339Nano), mj.Status.String()))
	}

	if err := t.saveJob(mj); err != nil {
		return 0, err
	}

	if err := t.tx.Set(mj.keyPayload(), payload); err != nil {
		return 0, err
	}

	// Create status index
	if err := t.tx.Set(mj.keyStatusIndex(mj.Status), []byte{}); err != nil {
		return 0, err
	}

	return newID, nil
}

func (t *Transaction) Update(ctx context.Context, jids []jobq.ID, updater func(job *jobq.JobInfo) error) error {
	for _, jid := range jids {

		mj, err := t.readJob(jid)
		if err != nil {
			fmt.Printf("$$$ job not found\n")
			return err
		}

		after := mj.ToDomain()
		if err := updater(after); err != nil {
			return err
		}

		save := false

		// Status changed ?
		if v := after.Status; v != mj.Status {
			if err := t.tx.Delete(mj.keyStatusIndex(mj.Status)); err != nil {
				fmt.Printf("delete key: %s: %v\n", mj.keyStatusIndex(mj.Status), err)
				return err
			}
			if err := t.tx.Set(mj.keyStatusIndex(v), []byte{}); err != nil {
				return err
			}
			mj.Status = v
			save = true
		}

		// DateTerminated changed ?
		if v := after.DateTerminated; v != mj.DateTerminated {
			mj.DateTerminated = v
			save = true
		}

		// DatesReserved changed ?
		if len(after.DatesReserved) != len(mj.DatesReserved) {
			mj.DatesReserved = after.DatesReserved
			save = true
		}

		// RetryCount changed ?
		if v := after.RetryCount; v != mj.RetryCount {
			mj.RetryCount = v
			save = true
		}

		// Logs changed ?
		if len(after.Logs) != len(mj.Logs) {
			mj.Logs = after.Logs
			save = true
		}

		if save {
			if err := t.saveJob(mj); err != nil {
				return err
			}
		}

	}

	return nil
}

func (t *Transaction) Delete(ctx context.Context, jids []jobq.ID) error {
	for _, jid := range jids {
		mj, err := t.readJob(jid)
		if err != nil {
			if errors.Is(repo.ErrJobNotFound, err) {
				continue
			}

			return err
		}

		if err := t.tx.Delete(mj.keyStatusIndex(mj.Status)); err != nil {
			return err
		}
		if err := t.tx.Delete(mj.keyPayload()); err != nil {
			return err
		}
		if err := t.tx.Delete(mj.keyJob()); err != nil {
			return err
		}
	}

	return nil
}

func (t *Transaction) Logf(ctx context.Context, jid jobq.ID, format string, args ...any) error {
	if format == "" {
		return nil
	}

	mj, err := t.readJob(jid)
	if err != nil {
		return err
	}

	log := time.Now().Format(time.RFC3339Nano) + ": " + fmt.Sprintf(format, args...)
	mj.Logs = append(mj.Logs, log)

	return t.saveJob(mj)
}
