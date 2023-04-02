package service

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

type IJobQueueService interface {
	// Add a new job.
	Enqueue(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error)

	// Reserve up to <limit> jobs.
	Reserve(ctx context.Context, topic jobq.Topic, limit int) ([]IJobService, error)

	// Available returns the number of jobs ready to be reserved.
	Available(ctx context.Context, topic jobq.Topic) (int, error)

	// Delayed returns the number of delayed jobs.
	Delayed(ctx context.Context) (int, error)

	// Topics returns list of created topics
	Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error)

	// TopicStats returns some stats about topic.
	TopicStats(ctx context.Context, topic jobq.Topic) (repo.TopicStats, error)

	GetJobs(ctx context.Context, []jobq.ID) ([]*jobq.JobInfo, error)
	
	Close() error
}

type IJobService interface {
	ID() jobq.ID
	Unwrap() (*jobq.JobInfo, error)
	Payload() (jobq.Payload, error)
	Logf(format string, args ...any) error
	Done() error
	Retry(overrideBackoff time.Duration) error
	Cancel() error
}
