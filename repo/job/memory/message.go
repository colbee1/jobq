package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

func (adapter *Adapter) AppendMessage(ctx context.Context, jid jobq.JobID, msg string) error {
	now := time.Now()

	if job, found := adapter.jobs[jid]; found {
		job.Message += msg
		job.DateUpdated = now
	}

	return nil
}
