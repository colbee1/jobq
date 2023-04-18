package service

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/topic"
)

//-----------------------------------------------------------------------
//	Job Service Interface
//-----------------------------------------------------------------------

type IJobService interface {
	// Add a new job.
	Enqueue(ctx context.Context, topic jobq.Topic, w jobq.Weight, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error)

	// Reserve up to <limit> jobs until context expiration.
	Reserve(ctx context.Context, topic jobq.Topic, limit int) ([]IJob, error)

	// Get current job status
	Status(ctx context.Context, jid jobq.ID) ([]jobq.Status, error)

	// Reset resets a job.
	Reset(ctx context.Context, jids []jobq.ID) error

	// Topics returns list of active topics.
	Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error)

	// TopicStats returns stats about topic.
	TopicStats(ctx context.Context, topic jobq.Topic) (topic.Stats, error)

	Close() error
}

type IJob interface {
	ID() jobq.ID
	Unwrap() (*jobq.JobInfo, error)
	Payload() (jobq.Payload, error)
	Logf(format string, args ...any) error
	Done() error
	Retry(overrideBackoff time.Duration) error
	Cancel() error
}
