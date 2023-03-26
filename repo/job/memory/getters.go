package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (t *Transaction) GetStatus(ctx context.Context, jid jobq.ID) (jobq.Status, error) {
	a := t.a

	if job, found := a.jobs[jid]; found {
		return job.status, nil
	}

	return 0, jobq.ErrJobNotFound
}

func (t *Transaction) GetPriority(ctx context.Context, jid jobq.ID) (jobq.Priority, error) {
	a := t.a

	if job, found := a.jobs[jid]; found {
		return job.priority, nil
	}

	return 0, jobq.ErrJobNotFound
}

func (a *Adapter) Durable() bool {
	return false
}

func (t *Transaction) GetInfos(ctx context.Context, jids []jobq.ID) ([]*jobq.JobInfo, error) {
	a := t.a

	if len(jids) == 0 {
		return []*jobq.JobInfo{}, nil
	}

	result := make([]*jobq.JobInfo, 0, len(jids))
	for _, jid := range jids {
		if m, found := a.jobs[jid]; found {
			ji := &jobq.JobInfo{
				ID:             m.id,
				Topic:          m.topic,
				Priority:       m.priority,
				Status:         m.status,
				DateCreated:    m.info.DateCreated,
				DateTerminated: m.info.DateTerminated,
				DateReserved:   m.info.DateReserved,
				Retries:        m.info.Retries,
				Options:        m.options,
				Logs:           m.logs,
			}
			result = append(result, ji)
		}
	}

	return result, nil
}

func (t *Transaction) GetOptions(ctx context.Context, jid jobq.ID) (*jobq.JobOptions, error) {
	a := t.a

	if m, found := a.jobs[jid]; found {
		return &jobq.JobOptions{
			Name:            m.options.Name,
			Timeout:         m.options.Timeout,
			DelayedAt:       m.options.DelayedAt,
			MaxRetries:      m.options.MaxRetries,
			MinBackOff:      m.options.MinBackOff,
			MaxBackOff:      m.options.MaxBackOff,
			LogStatusChange: m.options.LogStatusChange,
		}, nil
	}
	return nil, jobq.ErrJobNotFound
}

func (t *Transaction) Logs(ctx context.Context, jid jobq.ID) ([]string, error) {
	a := t.a

	job, found := a.jobs[jid]
	if !found {
		return nil, jobq.ErrJobNotFound
	}

	return job.logs, nil
}

func (t *Transaction) ListByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error) {
	a := t.a

	if limit < 1 {
		limit = 1000
	}

	jobs := make([]jobq.ID, 0, limit)

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
