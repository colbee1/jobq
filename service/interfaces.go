package service

import (
	"context"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

type IJobQueueService interface {
	// Add a new job.
	Enqueue(ctx context.Context, topic jobq.Topic, pri jobq.Priority, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error)

	// Reserve up to <limit> jobs.
	Reserve(ctx context.Context, topic jobq.Topic, limit int) ([]jobq.IJob, error)

	// Available returns the number of jobs ready to be reserved.
	Available(ctx context.Context, topic jobq.Topic) (int, error)

	// Delayed returns the number of delayed jobs.
	Delayed(ctx context.Context) (int, error)

	// Topics returns list of created topics
	Topics(ctx context.Context, offset int, limit int) ([]jobq.Topic, error)

	// TopicStats returns some stats about topic.
	TopicStats(ctx context.Context, topic jobq.Topic) (repo.TopicStats, error)

	Close() error
}
