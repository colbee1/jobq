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

	return &Service{jobRepo: jobRepo, pqRepo: pqRepo}, nil
}

func (s *Service) Close() error {
	return nil
}
