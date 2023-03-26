package badger3

import (
	"github.com/colbee1/jobq/repo"
	"github.com/dgraph-io/badger/v3"
)

type Transaction struct {
	a            *Adapter
	tx           *badger.Txn
	needCommit   bool
	commitFailed bool
}

func (a *Adapter) NewTransaction() (repo.IJobRepositoryTransaction, error) {
	return &Transaction{
		a:  a,
		tx: a.db.NewTransaction(true),
	}, nil
}

func (t *Transaction) Commit() error {
	if t.a == nil {
		return repo.ErrInvalidTransaction
	}

	var err error
	if t.needCommit {
		err = t.tx.Commit()
		if err != nil {
			t.commitFailed = true
		}
	}

	return err
}

func (t *Transaction) Close() error {
	t.tx.Discard()
	t.a = nil

	return nil
}
