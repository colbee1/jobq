package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (adapter *Adapter) SetStatus(ctx context.Context, jids []jobq.JobID, state jobq.JobStatus) error {
	now := time.Now()

	for _, jid := range jids {
		if job, found := adapter.jobs[jid]; found {
			job.Status = state
			job.DateUpdated = now
		}
	}

	return nil
}

func (adapter *Adapter) GetStatus(ctx context.Context, jid jobq.JobID) (jobq.JobStatus, error) {
	if job, found := adapter.jobs[jid]; found {
		return job.Status, nil
	}

	return 0, jobq.ErrJobNotFound
}
