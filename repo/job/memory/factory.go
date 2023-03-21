package memory

import "github.com/colbee1/jobq"

func New() (*Adapter, error) {
	return &Adapter{
		jobs: make(map[jobq.JobID]*jobq.Job),
	}, nil
}

func (a *Adapter) Close() error {
	a.jobs = make(map[jobq.JobID]*jobq.Job)

	return nil
}
