package service

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
)

const DefaultTopic = "$def"

type JobServiceInterface interface {
	// Enqueues enqueues a new job
	Enqueue(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jo jobq.JobOptions) (*Job, error)

	// Reserve reserves up to <limit> ready job
	Reserve(ctx context.Context, topic jobq.JobTopic, limit int) ([]jobq.JobID, error)
}

type Service struct {
	jobRepo JobRepositoryInterface
	pqRepo  JobPriorityQueueRepositoryInterface
}

type Job struct {
	service  *Service
	jid      jobq.JobID
	topic    jobq.JobTopic
	priority jobq.JobPriority
}

func (j *Job) ID() jobq.JobID {
	return j.jid
}

func (j *Job) AppendMessage(ctx context.Context, message string) error {
	return j.service.jobRepo.AppendMessage(ctx, j.jid, message)
}

func (j *Job) Release(ctx context.Context, delayedAt time.Time) error {
	jid := j.jid
	s := j.service

	// Job status must be reserved
	//
	state, err := s.jobRepo.GetStatus(ctx, jid)
	if err != nil {
		return err
	}
	if state != jobq.JobStateReserved {
		return fmt.Errorf("%w: job must be in reserved state to be relased", jobq.ErrInvalidJobState)
	}

	// Update job status and push it back in ready queue
	//
	if err := s.jobRepo.SetStatus(ctx, []jobq.JobID{jid}, jobq.JobStateReady); err != nil {
		return err
	}
	if err := s.pqRepo.Push(ctx, j.topic, j.priority, jid, delayedAt); err != nil {
		return err
	}

	return nil
}

func (j *Job) Done(ctx context.Context) error {
	jid := j.jid
	s := j.service

	// Job status must be reserved
	//
	state, err := s.jobRepo.GetStatus(ctx, jid)
	if err != nil {
		return err
	}
	if state != jobq.JobStateReserved {
		return fmt.Errorf("%w: job must be in reserved state to be relased", jobq.ErrInvalidJobState)
	}

	// Update job status
	//
	if err := s.jobRepo.SetStatus(ctx, []jobq.JobID{jid}, jobq.JobStateDone); err != nil {
		return err
	}

	return nil
}

func (j *Job) Error(ctx context.Context) error {
	jid := j.jid
	s := j.service

	// Job status must be reserved
	//
	state, err := s.jobRepo.GetStatus(ctx, jid)
	if err != nil {
		return err
	}
	if state != jobq.JobStateReserved {
		return fmt.Errorf("%w: job must be in reserved state to be relased", jobq.ErrInvalidJobState)
	}

	// Update job status
	//
	if err := s.jobRepo.SetStatus(ctx, []jobq.JobID{jid}, jobq.JobStateFail); err != nil {
		return err
	}

	return nil
}
