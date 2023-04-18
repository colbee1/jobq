package badger3

import (
	"context"

	"github.com/colbee1/jobq"
	repo "github.com/colbee1/jobq/repo/job"
	"github.com/dgraph-io/badger/v3"
)

type Adapter struct {
	db       *badger.DB
	jobIdSeq *badger.Sequence
}

type RepositoryOptions struct {
	SyncWrite bool
	DropDB    bool // Only for testing
}

func (a *Adapter) Durable() bool {
	return true
}

func (a *Adapter) Status(ctx context.Context, jid jobq.ID) (jobq.Status, error) {
	// TODO: Create an index to store job status
	tx, err := a.NewTransaction()
	if err != nil {
		return jobq.JobStatusUndefined, err
	}
	defer tx.Close()

	mjs, err := tx.Read(ctx, []jobq.ID{jid})
	if err != nil {
		return jobq.JobStatusUndefined, err
	}

	if len(mjs) == 0 {
		return jobq.JobStatusUndefined, repo.ErrJobNotFound
	}

	return mjs[0].Status, nil
}
