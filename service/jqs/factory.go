package jqs

import (
	"github.com/colbee1/assertor"
	"github.com/colbee1/jobq/repo/job"
	"github.com/colbee1/jobq/repo/topic"
)

func New(jobRepo job.IJobRepository, pqRepo topic.ITopicRepository) (*Service, error) {
	v := assertor.New()
	v.Assert(jobRepo != nil, "job repository is missing")
	v.Assert(pqRepo != nil, "job priority queue repository is missing")
	if err := v.Validate(); err != nil {
		return nil, err
	}

	s := &Service{
		jobRepo:       jobRepo,
		topicRepo:     pqRepo,
		stopScheduler: make(chan struct{}, 1),
	}

	s.wg.Add(1)
	go s.scheduler(s.stopScheduler)

	return s, nil
}

func (s *Service) Close() error {
	close(s.stopScheduler)
	s.wg.Wait()

	return nil
}
