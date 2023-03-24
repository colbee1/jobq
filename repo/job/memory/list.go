package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (t *Transaction) ListByStatus(ctx context.Context, status jobq.JobStatus, offset int, limit int) ([]jobq.JobID, error) {
	a := t.a

	if limit < 1 {
		limit = 1000
	}

	jobs := make([]jobq.JobID, 0, limit)

	a.jobsLock.RLock()
	defer a.jobsLock.RUnlock()

	for _, job := range a.jobs {
		if offset > 0 {
			offset--
			continue
		}

		if job.status == status {
			jobs = append(jobs, job.id)
		}

		if len(jobs) == limit {
			break
		}
	}

	return jobs, nil
}
