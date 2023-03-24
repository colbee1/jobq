package jobq

import (
	"context"
	"time"
)

type IJobPriorityQueueRepository interface {
	// Push adds a new jobid in job queue.
	Push(ctx context.Context, topic JobTopic, pri JobPriority, jid JobID, delayedAt time.Time) (JobStatus, error)

	// Pop returns up to <limit> jobs ordered by priority.
	Pop(ctx context.Context, topic JobTopic, limit int) ([]JobID, error)

	// Len returns the number of jobs in priority queue.
	Len(ctx context.Context, topic JobTopic) (int, error)

	// Durable returns true when repository can survive to an application crash.
	Durable() bool

	// Topics returns list of created topics
	Topics(ctx context.Context, offset int, limit int) ([]JobTopic, error)

	TopicStats(ctx context.Context, topic JobTopic) (TopicStats, error)

	Close() error
}

type TopicStats struct {
	DateCreated    time.Time
	DateLastPush   time.Time
	PushTotalCount int64
	MaxQueueLen    int64
}
