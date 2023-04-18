package badger3

import (
	repo "github.com/colbee1/jobq/repo/job"
	"github.com/dgraph-io/badger/v3"
)

func New(dbPath string, options RepositoryOptions) (*Adapter, error) {
	dbOpts := badger.DefaultOptions(dbPath).
		WithSyncWrites(options.SyncWrite).
		WithLogger(nil)
	db, err := badger.Open(dbOpts)
	if err != nil {
		return nil, err
	}

	seq, err := db.GetSequence(prefixKeyIdSequence, 1000)
	if err != nil {
		db.Close()
		return nil, err
	}

	a := &Adapter{
		db:       db,
		jobIdSeq: seq,
	}

	if options.DropDB {
		err := a.db.DropAll()
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	return a, nil
}

func (a *Adapter) Close() error {
	if a.jobIdSeq != nil {
		a.jobIdSeq.Release()
	}

	if a.db != nil && !a.db.IsClosed() {
		a.db.RunValueLogGC(0.5) // TODO: Remember to move this operation in regular time interval.
		a.db.Flatten(8)

		return a.db.Close()
	}

	return nil
}

func (a *Adapter) NewTransaction() (repo.IJobRepositoryTransaction, error) {
	return &Transaction{
		a:  a,
		tx: a.db.NewTransaction(true),
	}, nil
}
