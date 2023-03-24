package jobq

import (
	"context"
)

type IJobRepository interface {
	NewTransaction() (IJobRepositoryTransaction, error)

	// Durable returns true when repository can survive to an application crash.
	Durable() bool

	Close() error
}

type IJobRepositoryTransaction interface {
	// Creates creates a new job and returns it's uinque id
	Create(ctx context.Context, topic JobTopic, pri JobPriority, jo JobOptions, payload JobPayload) (JobID, error)

	SetStatus(ctx context.Context, jids []JobID, state JobStatus) error
	GetStatus(ctx context.Context, jid JobID) (JobStatus, error)

	// Log adds a log message in job logs. Date (rfc3339) is prepended.
	Log(ctx context.Context, jid JobID, message string) error

	// Logs returns all recorded logs.
	Logs(ctx context.Context, jid JobID) ([]string, error)

	ListByStatus(ctx context.Context, status JobStatus, offset int, limit int) ([]JobID, error)

	GetInfos(ctx context.Context, jids []JobID) ([]*JobInfo, error)
	GetPriority(ctx context.Context, jid JobID) (JobPriority, error)
	GetOptions(ctx context.Context, jid JobID) (*JobOptions, error)
	SetOptions(ctx context.Context, jid JobID, jo *JobOptions) error

	// GetPayload(ctx context.Context, jid JobID) (JobPayload, error)

	// Delete deletes done and canceled job by their ID.
	// DeleteByID(ctx context.Context, jids []JobID) error

	// Delete deletes done and canceled job by thei age.
	// DeleteByAge(ctx context.Context, age time.Duration) error

	Commit() error
	Close() error
}
