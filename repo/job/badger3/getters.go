package badger3

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

func (t *Transaction) readJob(jid jobq.ID) (*modelJob, error) {
	mj := &modelJob{ID: jid}
	itm, err := t.tx.Get(mj.keyJob())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("readJob(): %w: id=%d, key=%s", repo.ErrJobNotFound, jid, mj.keyJob())
		}

		return nil, err
	}
	data, err := itm.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	if err := mj.Decode(data); err != nil {
		return nil, err
	}

	return mj, nil
}

func (t *Transaction) Read(ctx context.Context, jids []jobq.ID) ([]*jobq.JobInfo, error) {
	if len(jids) == 0 {
		return []*jobq.JobInfo{}, nil
	}

	jobs := make([]*jobq.JobInfo, 0, len(jids))
	for _, jid := range jids {
		mj, err := t.readJob(jid)
		if err != nil {
			if errors.Is(err, repo.ErrJobNotFound) {
				continue
			}

			return nil, err
		}

		jobs = append(jobs, mj.ToDomain())
	}

	return jobs, nil
}

func (t *Transaction) ReadPayload(ctx context.Context, jid jobq.ID) (jobq.Payload, error) {
	m := &modelJob{ID: jid}
	itm, err := t.tx.Get(m.keyPayload())
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

	return jobq.Payload(data), nil
}

func (t *Transaction) FindByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error) {
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
