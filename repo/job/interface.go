package job

import (
	"context"

	"github.com/colbee1/jobq"
)

type IJobRepository interface {
	// NewTransaction Creates a new transaction.
	NewTransaction() (IJobRepositoryTransaction, error)

	// Durable returns true when repository can survive to an application crash.
	Durable() bool

	Close() error
}

type IJobRepositoryTransaction interface {
	// Creates creates a new job and returns it's unique id
	Create(ctx context.Context, topic jobq.Topic, w jobq.Weight, jo jobq.JobOptions, payload jobq.Payload) (jobq.ID, error)

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
