package badger3

import (
	"github.com/dgraph-io/badger/v3"
)

func New(dbPath string, options Options) (*Adapter, error) {
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

	if options.DropAll {
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

	return a.db.Close()
}
