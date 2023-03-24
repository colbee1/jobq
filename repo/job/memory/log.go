package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (t *Transaction) Log(ctx context.Context, jid jobq.JobID, msg string) error {
	a := t.a

	if job, found := a.jobs[jid]; !found {
		return jobq.ErrJobNotFound
	} else {
		job.log(msg)
		t.needCommit = true
	}

	return nil
}

func (t *Transaction) Logs(ctx context.Context, jid jobq.JobID) ([]string, error) {
	a := t.a

	job, found := a.jobs[jid]
	if !found {
		return nil, jobq.ErrJobNotFound
	}

	return job.logs, nil
}
