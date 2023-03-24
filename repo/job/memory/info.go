package memory

import (
	"context"

	"github.com/colbee1/jobq"
)

func (t *Transaction) GetInfos(ctx context.Context, jids []jobq.JobID) ([]*jobq.JobInfo, error) {
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

func (t *Transaction) GetOptions(ctx context.Context, jid jobq.JobID) (*jobq.JobOptions, error) {
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

func (t *Transaction) SetOptions(ctx context.Context, jid jobq.JobID, jo *jobq.JobOptions) error {
	a := t.a

	if m, found := a.jobs[jid]; found {
		m.options = *jo

		return nil
	}

	return jobq.ErrJobNotFound
}
