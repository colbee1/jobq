package topic

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

type ITopicRepository interface {
	// Push adds a new jobid in desired topic.
	Push(ctx context.Context, topic jobq.Topic, w jobq.Weight, jid jobq.ID, delayedAt time.Time) (jobq.Status, error)

	// Pop returns up to <limit> jobs from desired topic.
	PopTopic(ctx context.Context, topic jobq.Topic, limit int) ([]jobq.ID, error)

	// PopDelayed returns up to <limit> jobs from delayed queue.
	PopDelayed(ctx context.Context, limit int) ([]jobq.ID, error)

	// Topics returns list of active topics.
	// Stats collector should be enabled.
	Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error)

	// TopicStats returns some stats on topic.
	// Stats collector should be enabled.
	// Use "delayed" as topic to get stats from delayed queue.
	TopicStats(ctx context.Context, topic jobq.Topic) (Stats, error)

	// Durable returns true when repository is persistant.
	Durable() bool

	Close() error
}

type Stats struct {
	PushDateFirst  time.Time
	PushDateLast   time.Time
	PushTotalCount int64
	PopDateFirst   time.Time
	PopDateLast    time.Time
	PopTotalCount  int64
}
