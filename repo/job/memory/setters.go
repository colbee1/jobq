package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

// Create creates a new job and returns it's ID
func (t *Transaction) Create(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
	if topic == "" {
		return 0, jobq.ErrTopicIsMissing
	}

	newID := jobq.ID(t.a.jobSequence.Add(1))

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
		Payload: payload,
	}

	if mj.Options.LogStatusChange {
		mj.Logs = append(mj.Logs, fmt.Sprintf("%s: Job status changed to: %s",
			time.Now().Format(time.RFC3339Nano), mj.Status.String()))
	}

	t.a.jobs[newID] = mj
	t.needCommit = true

	return newID, nil
}

func (t *Transaction) Update(ctx context.Context, jids []jobq.ID, updater func(job *jobq.JobInfo) error) error {
	t.a.jobsLock.Lock()
	defer t.a.jobsLock.Unlock()

	for _, jid := range jids {

		mj, found := t.a.jobs[jid]
		if !found {
			return repo.ErrJobNotFound
		}

		after := mj.ToDomain()
		if err := updater(after); err != nil {
			return err
		}

		// Status changed ?
		if v := after.Status; v != mj.Status {
			mj.Status = v
		}

		// DateTerminated changed ?
		if v := after.DateTerminated; v != mj.DateTerminated {
			mj.DateTerminated = v
		}

		// DatesReserved changed ?
		if len(after.DatesReserved) != len(mj.DatesReserved) {
			mj.DatesReserved = after.DatesReserved
		}

		// RetryCount changed ?
		if v := after.RetryCount; v != mj.RetryCount {
			mj.RetryCount = v
		}

		// Logs changed ?
		if len(after.Logs) != len(mj.Logs) {
			mj.Logs = after.Logs
		}

	}

	return nil
}

func (t *Transaction) Delete(ctx context.Context, jids []jobq.ID) error {
	t.a.jobsLock.Lock()
	defer t.a.jobsLock.Unlock()

	for _, jid := range jids {
		delete(t.a.jobs, jid)
	}

	return nil
}

func (t *Transaction) Logf(ctx context.Context, jid jobq.ID, format string, args ...any) error {
	if format == "" {
		return nil
	}

	t.a.jobsLock.Lock()
	defer t.a.jobsLock.Unlock()

	mj, found := t.a.jobs[jid]
	if !found {
		return repo.ErrJobNotFound
	}

	log := time.Now().Format(time.RFC3339Nano) + ": " + fmt.Sprintf(format, args...)
	mj.Logs = append(mj.Logs, log)

	return nil
}
