package service

import "github.com/colbee1/assertor"

func New(jobRepo JobRepositoryInterface, pqRepo JobPriorityQueueRepositoryInterface) (*Service, error) {
	v := assertor.New()
	v.Assert(jobRepo != nil, "job repository is missing")
	v.Assert(pqRepo != nil, "job priority queue repository is missing")
	if err := v.Validate(); err != nil {
		return nil, err
	}

	return &Service{jobRepo: jobRepo, pqRepo: pqRepo}, nil
}
