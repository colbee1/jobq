package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (adapter *Adapter) Find(ctx context.Context, jids []jobq.JobID) ([]*jobq.Job, error) {
	if len(jids) == 0 {
		return []*jobq.Job{}, nil
	}

	jobs := make([]*jobq.Job, 0, len(jids))
	adapter.jobsLock.RLock()
	defer adapter.jobsLock.RUnlock()
	for _, jid := range jids {
		if job, found := adapter.jobs[jid]; found {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (adapter *Adapter) FindByStatus(ctx context.Context, status jobq.JobStatus, offset int, limit int) ([]*jobq.Job, error) {
	if limit < 1 {
		limit = 1000
	}

	jobs := make([]*jobq.Job, 0, limit)

	adapter.jobsLock.RLock()
	defer adapter.jobsLock.RUnlock()

	for _, job := range adapter.jobs {
		if offset > 0 {
			offset--
			continue
		}

		if job.Status == status {
			jobs = append(jobs, job)
		}

		if len(jobs) == limit {
			break
		}
	}

	return jobs, nil
}
