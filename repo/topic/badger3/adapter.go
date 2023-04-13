package badger3

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/colbee1/assertor"
	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/topic"
	"github.com/dgraph-io/badger/v3"
)

type Adapter struct {
	db             *badger.DB
	statsCollector *topic.StatsCollector
}

type RepositoryOptions struct {
	SyncWrite      bool
	StatsCollector bool
	DropDB         bool
}

func (a *Adapter) createOrGetTopic(tx *badger.Txn, topic jobq.Topic) error {
	key := newKeyTopicInfo(topic)
	_, err := tx.Get(key.Info())
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	mti := &modelTopicInfo{
		DateCreated: time.Now(),
	}
	data, err := mti.Encode()
	if err != nil {
		return err
	}

	if err := tx.Set(key.Info(), data); err != nil {
		return err
	}

	return nil
}

func (a *Adapter) Push(ctx context.Context, tp jobq.Topic, w jobq.Weight, jid jobq.ID, delayedAt time.Time) (jobq.Status, error) {
	v := assertor.New()
	v.Assert(ctx != nil, "context is missing")
	v.Assert(tp != "", "topic is missing")
	v.Assert(jid > 0, "invalid job ID")
	if err := v.Validate(); err != nil {
		return jobq.JobStatusUndefined, err
	}

	tx := a.db.NewTransaction(true)
	defer tx.Discard()

	if err := a.createOrGetTopic(tx, tp); err != nil {
		return jobq.JobStatusUndefined, err
	}

	key := newKeyTopicItem(tp, w, jid).Item()
	status := jobq.JobStatusReady
	if time.Until(delayedAt) > time.Second {
		key = newKeyDelayedItem(delayedAt, jid).Item()
		status = jobq.JobStatusDelayed
	}
	if err := tx.Set(key, []byte{}); err != nil {
		return jobq.JobStatusUndefined, err
	}

	if a.statsCollector != nil {
		t := tp
		if status == jobq.JobStatusDelayed {
			t = "delayed"
		}
		a.statsCollector.Collector <- topic.StatMetric{Topic: t, Metric: topic.MetricPush, Value: 1}
	}

	return status, tx.Commit()
}

func (a *Adapter) PopTopic(ctx context.Context, tp jobq.Topic, limit int) ([]jobq.ID, error) {
	v := assertor.New()
	v.Assert(ctx != nil, "context is missing")
	v.Assert(tp != "", "topic is missing")
	v.Assert(limit > 0, "limit should be a non zero positive number")
	if err := v.Validate(); err != nil {
		return nil, err
	}

	jids := make([]jobq.ID, 0, limit)
	prefix := newKeyTopicItem(tp, 0, 0).Prefix()

	tx := a.db.NewTransaction(true)
	defer tx.Discard()

	poper := func() error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // access to the LSM-tree only
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			rawKey := it.Item().KeyCopy(nil) // Use KeyCopy because tx.Delete() keep reference on it.
			key := new(keyTopicItem)
			if err := key.Decode(rawKey); err != nil {
				return err
			}
			jids = append(jids, key.JobID)

			if err := tx.Delete(rawKey); err != nil {
				return err
			}
			fmt.Printf("deleting key: %v\n", string(rawKey))

			if a.statsCollector != nil {
				a.statsCollector.Collector <- topic.StatMetric{Topic: key.Topic, Metric: topic.MetricPop, Value: 1}
			}

			limit--
			if limit == 0 {
				break
			}

			if err := ctx.Err(); err != nil {
				return err
			}
		}

		return nil
	}

	if err := poper(); err == nil {
		return jids, tx.Commit()
	} else {
		return nil, err
	}
}

func (a *Adapter) PopDelayed(ctx context.Context, limit int) ([]jobq.ID, error) {
	v := assertor.New()
	v.Assert(ctx != nil, "context is missing")
	v.Assert(limit > 0, "limit should be a non zero positive number")
	if err := v.Validate(); err != nil {
		return nil, err
	}

	jids := make([]jobq.ID, 0, limit)
	prefix := newKeyDelayedItem(time.Time{}, 0).Prefix()

	tx := a.db.NewTransaction(true)
	defer tx.Discard()

	poper := func() error {
		now := time.Now()
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // access to the LSM-tree only
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			rawKey := it.Item().KeyCopy(nil) // Use KeyCopy because tx.Delete() keep reference on it.
			key := new(keyDelayedItem)
			if err := key.Decode(rawKey); err != nil {
				return err
			}

			if now.After(key.At) {
				jids = append(jids, key.JobID)
			} else {
				break
			}

			if err := tx.Delete(rawKey); err != nil {
				return err
			}

			if a.statsCollector != nil {
				// TODO: Fix static topic
				a.statsCollector.Collector <- topic.StatMetric{Topic: "delayed", Metric: topic.MetricPop, Value: 1}
			}

			limit--
			if limit == 0 {
				break
			}

			if err := ctx.Err(); err != nil {
				return err
			}
		}

		return nil
	}

	if err := poper(); err == nil {
		return jids, tx.Commit()
	} else {
		return nil, err
	}
}

func (a *Adapter) Durable() bool {
	return true
}

func (a *Adapter) Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error) {
	if a.statsCollector == nil {
		return nil, topic.ErrStatsCollectorDisabled
	}

	return a.statsCollector.Topics(int64(offset), int64(limit))
}

func (a *Adapter) TopicStats(ctx context.Context, tp jobq.Topic) (topic.Stats, error) {
	if a.statsCollector == nil {
		return topic.Stats{}, topic.ErrStatsCollectorDisabled
	}

	return a.statsCollector.TopicStats(tp)
}
