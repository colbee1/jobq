package jqs

import (
	"github.com/colbee1/assertor"
	"github.com/colbee1/jobq/repo"
)

func New(jobRepo repo.IJobRepository, pqRepo repo.IJobPriorityQueueRepository) (*Service, error) {
	v := assertor.New()
	v.Assert(jobRepo != nil, "job repository is missing")
	v.Assert(pqRepo != nil, "job priority queue repository is missing")
	if err := v.Validate(); err != nil {
		return nil, err
	}

	s := &Service{
		jobRepo:   jobRepo,
		pqRepo:    pqRepo,
		exitSched: make(chan struct{}, 1),
	}

	s.wg.Add(1)
	go s.scheduleDelayedJobs(s.exitSched)

	return s, nil
}

func (s *Service) Close() error {
	close(s.exitSched)
	s.wg.Wait()

	return nil
}
