package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (t *Transaction) SetStatus(ctx context.Context, jids []jobq.ID, status jobq.Status) error {
	a := t.a
	now := time.Now()

	for _, jid := range jids {
		if job, found := a.jobs[jid]; found {
			job.status = status
			t.needCommit = true

			switch status {
			case jobq.JobStatusDone, jobq.JobStatusCanceled:
				job.info.DateTerminated = now
			case jobq.JobStatusReserved:
				job.info.DateReserved = append(job.info.DateReserved, now)
			}

			if job.options.LogStatusChange {
				job.log("Status change to: " + status.String())
			}
		}
	}

	return nil
}

func (t *Transaction) Log(ctx context.Context, jid jobq.ID, msg string) error {
	a := t.a

	if job, found := a.jobs[jid]; !found {
		return jobq.ErrJobNotFound
	} else {
		job.log(msg)
		t.needCommit = true
	}

	return nil
}

func (t *Transaction) SetOptions(ctx context.Context, jid jobq.ID, jo *jobq.JobOptions) error {
	a := t.a

	if m, found := a.jobs[jid]; found {
		m.options = *jo

		return nil
	}

	return jobq.ErrJobNotFound
}

// Create creates a new job and returns it's ID
func (t *Transaction) Create(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
	a := t.a

	if topic == "" {
		return 0, jobq.ErrTopicIsMissing
	}

	jid := jobq.ID(a.jobSequence.Add(1))
	model := &modelJob{
		id:       jid,
		topic:    topic,
		priority: pri,
		status:   jobq.JobStatusCreated,
		payload:  payload,
		options:  jo,
		info:     jobq.JobInfo{DateCreated: time.Now()},
		logs:     []string{},
	}

	if model.options.LogStatusChange {
		model.log("job created")
	}

	a.jobs[jid] = model
	t.needCommit = true

	return jid, nil
}
