package memory

import (
	repo "github.com/colbee1/jobq/repo/job"
)

// memory adapter doesn't support transaction

type Transaction struct {
	a          *Adapter
	needCommit bool
}

func (t *Transaction) Commit() error {
	if t.a == nil {
		return repo.ErrInvalidTransaction
	}

	if !t.needCommit {
		return nil
	}

	return nil
}

func (t *Transaction) Close() error {
	t.a = nil

	return nil
}
