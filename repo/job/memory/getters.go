package memory

import (
	"context"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

func (a *Adapter) Durable() bool {
	return false
}

func (t *Transaction) Read(ctx context.Context, jids []jobq.ID) ([]*jobq.JobInfo, error) {
	a := t.a

	if len(jids) == 0 {
		return []*jobq.JobInfo{}, nil
	}

	result := make([]*jobq.JobInfo, 0, len(jids))
	for _, jid := range jids {
		if mj, found := a.jobs[jid]; found {
			result = append(result, mj.ToDomain())
		}
	}

	return result, nil
}

func (t *Transaction) ReadPayload(ctx context.Context, jid jobq.ID) (jobq.Payload, error) {
	if job, found := t.a.jobs[jid]; found {
		return job.Payload, nil
	}

	return nil, repo.ErrJobNotFound
}

func (t *Transaction) FindByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error) {
	a := t.a

	if limit < 1 {
		limit = 1000
	}

	jobs := make([]jobq.ID, 0, limit)

	a.jobsLock.RLock()
	defer a.jobsLock.RUnlock()

	for _, mj := range a.jobs {
		if offset > 0 {
			offset--
			continue
		}

		if mj.Status == status {
			jobs = append(jobs, mj.ID)
		}

		if len(jobs) == limit {
			break
		}
	}

	return jobs, nil
}
