package badger3

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/colbee1/jobq"
	repo "github.com/colbee1/jobq/repo/job"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

type Transaction struct {
	a            *Adapter
	tx           *badger.Txn
	needCommit   bool
	commitFailed bool
}

func (t *Transaction) Create(ctx context.Context, topic jobq.Topic, w jobq.Weight, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error) {
	if topic == "" {
		return 0, jobq.ErrTopicIsMissing
	}

	newID := jobq.ID(0)
	for newID == 0 {
		if id, err := t.a.jobIdSeq.Next(); err != nil {
			return 0, err
		} else {
			newID = jobq.ID(id)
		}
	}

	mj := &modelJob{
		ID:          newID,
		Topic:       topic,
		Weight:      w,
		Status:      jobq.JobStatusCreated,
		DateCreated: time.Now(),
		Options: modelJobOptions{
			Name:            jo.Name,
			Timeout:         jo.Timeout,
			DelayedAt:       jo.DelayedAt,
			MaxRetries:      jo.MaxRetries,
			MinBackOff:      jo.MinBackOff,
			MaxBackOff:      jo.MaxBackOff,
			LogStatusChange: jo.LogStatusChange,
		},
	}

	if mj.Options.LogStatusChange {
		mj.Logs = append(mj.Logs, fmt.Sprintf("%s: Job status changed to: %s",
			time.Now().Format(time.RFC3339Nano), mj.Status.String()))
	}

	if err := t.saveJob(mj); err != nil {
		return 0, err
	}

	if err := t.tx.Set(mj.keyPayload(), payload); err != nil {
		return 0, err
	}

	// Create status index
	if err := t.tx.Set(mj.keyStatusIndex(mj.Status), []byte{}); err != nil {
		return 0, err
	}

	return newID, nil
}

func (t *Transaction) Update(ctx context.Context, jids []jobq.ID, updater func(job *jobq.JobInfo) error) error {
	for _, jid := range jids {

		mj, err := t.readJob(jid)
		if err != nil {
			return err
		}

		after := mj.ToDomain()
		if err := updater(after); err != nil {
			return err
		}

		save := false

		// Status changed ?
		if v := after.Status; v != mj.Status {
			if err := t.tx.Delete(mj.keyStatusIndex(mj.Status)); err != nil {
				return err
			}
			if err := t.tx.Set(mj.keyStatusIndex(v), []byte{}); err != nil {
				return err
			}
			mj.Status = v
			save = true
		}

		// DateTerminated changed ?
		if v := after.DateTerminated; v != mj.DateTerminated {
			mj.DateTerminated = v
			save = true
		}

		// DatesReserved changed ?
		if len(after.DatesReserved) != len(mj.DatesReserved) {
			mj.DatesReserved = after.DatesReserved
			save = true
		}

		// RetryCount changed ?
		if v := after.RetryCount; v != mj.RetryCount {
			mj.RetryCount = v
			save = true
		}

		// ResetCount changed ?
		if v := after.ResetCount; v != mj.ResetCount {
			mj.ResetCount = v
			save = true
		}

		// Logs changed ?
		if len(after.Logs) != len(mj.Logs) {
			mj.Logs = after.Logs
			save = true
		}

		if save {
			if err := t.saveJob(mj); err != nil {
				return err
			}
		}

	}

	return nil
}

func (t *Transaction) Delete(ctx context.Context, jids []jobq.ID) error {
	// Create array of keys that can be referenced over all the transaction by Delete()
	keysToDelete := make([][]byte, 0, len(jids)*3)
	for _, jid := range jids {
		mj, err := t.readJob(jid)
		if err != nil {
			if errors.Is(repo.ErrJobNotFound, err) {
				continue
			}

			return err
		}
		keysToDelete = append(keysToDelete,
			mj.keyStatusIndex(mj.Status),
			mj.keyPayload(),
			mj.keyJob(),
		)
	}

	for _, keys := range keysToDelete {
		if err := t.tx.Delete(keys); err != nil {
			return err
		}
	}

	return nil
}

func (t *Transaction) Logf(ctx context.Context, jid jobq.ID, format string, args ...any) error {
	if format == "" {
		return nil
	}

	mj, err := t.readJob(jid)
	if err != nil {
		return err
	}

	log := time.Now().Format(time.RFC3339Nano) + ": " + fmt.Sprintf(format, args...)
	mj.Logs = append(mj.Logs, log)

	return t.saveJob(mj)
}

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

func (t *Transaction) saveJob(mj *modelJob) error {
	data, err := mj.Encode()
	if err != nil {
		return err
	}

	if err := t.tx.Set(mj.keyJob(), data); err != nil {
		return err
	}
	t.needCommit = true

	return nil
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

func (t *Transaction) Commit() error {
	if t.a == nil {
		return repo.ErrInvalidTransaction
	}

	var err error
	if t.needCommit {
		err = t.tx.Commit()
		if err != nil {
			t.commitFailed = true
		}
	}

	return err
}

func (t *Transaction) Close() error {
	t.tx.Discard()
	t.a = nil

	return nil
}
