package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (t *Transaction) SetStatus(ctx context.Context, jids []jobq.JobID, status jobq.JobStatus) error {
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

func (t *Transaction) GetStatus(ctx context.Context, jid jobq.JobID) (jobq.JobStatus, error) {
	a := t.a

	if job, found := a.jobs[jid]; found {
		return job.status, nil
	}

	return 0, jobq.ErrJobNotFound
}

func (t *Transaction) GetPriority(ctx context.Context, jid jobq.JobID) (jobq.JobPriority, error) {
	a := t.a

	if job, found := a.jobs[jid]; found {
		return job.priority, nil
	}

	return 0, jobq.ErrJobNotFound
}
