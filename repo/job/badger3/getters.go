package badger3

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
	"github.com/dgraph-io/badger/v3"
)

func (t *Transaction) GetStatus(ctx context.Context, jid jobq.ID) (jobq.Status, error) {
	status := jobq.JobStatusUndefined

	m := &modelJob{ID: jid}
	itm, err := t.tx.Get(m.keyStatus())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return status, repo.ErrJobNotFound
		}

		return status, err
	}

	data, err := itm.ValueCopy(nil)
	if err != nil {
		return status, err
	}

	return jobq.NewStatusFromString(string(data)), nil
}

func (t *Transaction) Logs(ctx context.Context, jid jobq.ID) ([]string, error) {
	tx := t.tx
	m := &modelJob{ID: jid}
	itm, err := tx.Get(m.keyLogs())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return []string{}, nil
		}

		return nil, err
	}

	data, err := itm.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	logs := []string{}

	return logs, json.Unmarshal(data, &logs)
}

func (t *Transaction) ListByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error) {
	result := []jobq.ID{}
	if limit == 0 {
		return result, nil
	}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(string(prefixKeyStatusIndex) + status.String() + ":")
	it := t.tx.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		if offset > 0 {
			offset--
			continue
		}

		item := it.Item()
		k := item.Key()
		index := bytes.LastIndex(k, []byte{':'})

		id, err := strconv.ParseUint(string(k[index+1:]), 10, 64)
		if err != nil {
			return result, err
		}
		jid := jobq.ID(id)
		result = append(result, jid)

		if len(result) == limit {
			break
		}
	}

	return result, nil
}

func (t *Transaction) GetInfos(ctx context.Context, jids []jobq.ID) ([]*jobq.JobInfo, error) {
	infos := []*jobq.JobInfo{}
	if len(jids) == 0 {
		return infos, nil
	}

	for _, jid := range jids {
		m := &modelJob{ID: jid}
		itm, err := t.tx.Get(m.keyJob())
		if err == badger.ErrKeyNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		data, err := itm.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		if err := m.Decode(data); err != nil {
			return nil, err
		}

		status, err := t.GetStatus(ctx, jid)
		if err != nil {
			return nil, err
		}

		logs, err := t.Logs(ctx, jid)
		if err != nil {
			return nil, err
		}

		info := &jobq.JobInfo{
			ID:             m.ID,
			Topic:          m.Topic,
			Priority:       m.Priority,
			Status:         status,
			DateCreated:    m.DateCreated,
			DateTerminated: m.DateTerminated,
			DateReserved:   []time.Time{}, // TODO:
			Retries:        0,             // TODO:
			Options: jobq.JobOptions{
				Name:            m.Options.Name,
				Timeout:         m.Options.Timeout,
				DelayedAt:       m.Options.DelayedAt,
				MaxRetries:      m.Options.MaxRetries,
				MinBackOff:      m.Options.MinBackOff,
				MaxBackOff:      m.Options.MaxBackOff,
				LogStatusChange: m.Options.LogStatusChange,
			},
			Logs: logs,
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func (t *Transaction) GetPriority(ctx context.Context, jid jobq.ID) (jobq.Priority, error) {
	pri := jobq.Priority(0)

	m := &modelJob{ID: jid}
	itm, err := t.tx.Get(m.keyJob())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return pri, repo.ErrJobNotFound
		}

		return pri, err
	}

	data, err := itm.ValueCopy(nil)
	if err != nil {
		return pri, err
	}
	if err := m.Decode(data); err != nil {
		return pri, err
	}

	return m.Priority, nil
}

func (t *Transaction) GetOptions(ctx context.Context, jid jobq.ID) (*jobq.JobOptions, error) {
	m := &modelJob{ID: jid}
	itm, err := t.tx.Get(m.keyJob())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, repo.ErrJobNotFound
		}

		return nil, err
	}

	data, err := itm.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	if err := m.Decode(data); err != nil {
		return nil, err
	}

	return &jobq.JobOptions{
		Name:            m.Options.Name,
		Timeout:         m.Options.Timeout,
		DelayedAt:       m.Options.DelayedAt,
		MaxRetries:      m.Options.MaxRetries,
		MinBackOff:      m.Options.MinBackOff,
		MaxBackOff:      m.Options.MaxBackOff,
		LogStatusChange: m.Options.LogStatusChange,
	}, nil
}
