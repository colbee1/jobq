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
	// Creates creates a new job and returns it's unique id
	Create(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error)

	// Read returns all informations about one jobs.
	Read(ctx context.Context, jids []jobq.ID) ([]*jobq.JobInfo, error)

	// Payload returns job's payload.
	ReadPayload(ctx context.Context, jid jobq.ID) (jobq.Payload, error)

	// Update applies job mutations on jobs.
	Update(ctx context.Context, jids []jobq.ID, updater func(job *jobq.JobInfo) error) error

	// Delete removes jobs for repository.
	Delete(ctx context.Context, jids []jobq.ID) error

	// FindByStatus returns list of jobq.IDs for job in wanted status.
	FindByStatus(ctx context.Context, status jobq.Status, offset int, limit int) ([]jobq.ID, error)

	// Logf adds a formated log message in job logs. Date (rfc3339) is prepended.
	Logf(ctx context.Context, jid jobq.ID, format string, args ...any) error

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
