package jobq

import "context"

type IJobQueueService interface {
	// Add a new job.
	Enqueue(ctx context.Context, topic JobTopic, pri JobPriority, jo JobOptions, payload JobPayload) (JobID, error)

	// Reserve up to <limit> jobs.
	Reserve(ctx context.Context, topic JobTopic, limit int) ([]IJob, error)

	// Available returns the number of jobs ready to be reserved.
	Available(ctx context.Context, topic JobTopic) (int, error)

	// Delayed returns the number of delayed jobs.
	Delayed(ctx context.Context) (int, error)

	Close() error
}
