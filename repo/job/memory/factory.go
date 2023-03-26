package memory

import (
	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

func New() (*Adapter, error) {
	return &Adapter{
		jobs: make(map[jobq.ID]*modelJob),
	}, nil
}

func (a *Adapter) NewTransaction() (repo.IJobRepositoryTransaction, error) {
	return &Transaction{a: a}, nil
}

func (a *Adapter) Close() error {
	a.jobs = make(map[jobq.ID]*modelJob)

	return nil
}
