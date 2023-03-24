package memory

import (
	"github.com/colbee1/jobq"
)

func New() (*Adapter, error) {
	return &Adapter{
		jobs: make(map[jobq.JobID]*modelJob),
	}, nil
}

func (a *Adapter) NewTransaction() (jobq.IJobRepositoryTransaction, error) {
	return &Transaction{a: a}, nil
}
