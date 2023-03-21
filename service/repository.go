package service

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

type JobRepositoryInterface interface {
	Create(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jo jobq.JobOptions) (jobq.JobID, error)
	Find(ctx context.Context, jids []jobq.JobID) ([]*jobq.Job, error)
	FindByStatus(ctx context.Context, status jobq.JobStatus, offset int, limit int) ([]*jobq.Job, error)
	AppendMessage(ctx context.Context, jid jobq.JobID, message string) error
	SetStatus(ctx context.Context, jids []jobq.JobID, state jobq.JobStatus) error
	GetStatus(ctx context.Context, jid jobq.JobID) (jobq.JobStatus, error)

	// Durable returns true when repository can survive to an application crash
	Durable() bool

	Close() error
}

type JobPriorityQueueRepositoryInterface interface {
	// Push adds a new jobid in job queue
	Push(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jid jobq.JobID, delayedAt time.Time) error

	// Pop returns the first <limit> jobs ordered by priority
	Pop(ctx context.Context, topic jobq.JobTopic, limit int) ([]jobq.JobID, error)

	// Len returns the number of jobs in priority queue
	Len(ctx context.Context, topic jobq.JobTopic) (int, error)

	// Durable returns true when repository can survive to an application crash
	Durable() bool

	Close() error
}
