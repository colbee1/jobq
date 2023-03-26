package repo

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

//
// Job repository interfaces
//

type IJobRepository interface {
	NewTransaction() (IJobRepositoryTransaction, error)

	// Durable returns true when repository can survive to an application crash.
	Durable() bool

	Close() error
}

type IJobRepositoryTransaction interface {
	// Creates creates a new job and returns it's uinque id
	Create(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error)

	SetStatus(ctx context.Context, jids []jobq.ID, state jobq.Status) error
	GetStatus(ctx context.Context, jid jobq.ID) (jobq.Status, error)

	// Log adds a log message in job logs. Date (rfc3339) is prepended.
	Log(ctx context.Context, jid jobq.ID, message string) error

	// Logs returns all recorded logs.
	Logs(ctx context.Context, jid jobq.ID) ([]string, error)

	ListByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error)

	GetInfos(ctx context.Context, jids []jobq.ID) ([]*jobq.JobInfo, error)
	GetPriority(ctx context.Context, jid jobq.ID) (jobq.Priority, error)
	GetOptions(ctx context.Context, jid jobq.ID) (*jobq.JobOptions, error)
	SetOptions(ctx context.Context, jid jobq.ID, jo *jobq.JobOptions) error

	// GetPayload(ctx context.Context, jid JobID) (JobPayload, error)

	// Delete deletes done and canceled job by their ID.
	// DeleteByID(ctx context.Context, jids []JobID) error

	// Delete deletes done and canceled job by thei age.
	// DeleteByAge(ctx context.Context, age time.Duration) error

	Commit() error
	Close() error
}

//
// Priority Queue repository interfaces
//

type IJobPriorityQueueRepository interface {
	// Push adds a new jobid in job queue.
	Push(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jid jobq.ID, delayedAt time.Time) (jobq.Status, error)

	// Pop returns up to <limit> jobs ordered by priority.
	Pop(ctx context.Context, topic jobq.Topic, limit int) ([]jobq.ID, error)

	// Len returns the number of jobs in priority queue.
	Len(ctx context.Context, topic jobq.Topic) (int, error)

	// Durable returns true when repository can survive to an application crash.
	Durable() bool

	// Topics returns list of created topics.
	Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error)

	// TopicStats returns some simple stats on topic.
	TopicStats(ctx context.Context, topic jobq.Topic) (TopicStats, error)

	Close() error
}

type TopicStats struct {
	DateCreated          time.Time
	DateLastPush         time.Time
	PushTotalCount       int64
	MaxQueueLen          int64
	CurrentQueueCapacity int64
}
