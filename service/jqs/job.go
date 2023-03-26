package jqs

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
)

type Job struct {
	service *Service
	id      jobq.ID
	topic   jobq.Topic
}

func (j *Job) ID() jobq.ID {
	return j.id
}

func (j *Job) Topic() jobq.Topic {
	return j.topic
}

func (j *Job) Priority() (jobq.Priority, error) {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return 0, err
	}
	defer tx.Close()

	return tx.GetPriority(context.Background(), j.id)
}

func (j *Job) Status() (jobq.Status, error) {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return jobq.JobStatusCreated, err
	}
	defer tx.Close()

	return tx.GetStatus(context.Background(), j.id)
}

func (j *Job) Log(message string) error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	if err := tx.Log(context.Background(), j.id, message); err != nil {
		return err
	}

	return tx.Commit()
}

func (j *Job) Logs() ([]string, error) {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	return tx.Logs(context.Background(), j.id)
}

func (j *Job) Done(log string) error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	// Job status must be reserved
	//
	status, err := tx.GetStatus(context.Background(), j.ID())
	if err != nil {
		return err
	}
	if status != jobq.JobStatusReserved {
		return fmt.Errorf("%w: job must be in reserved state to be done", jobq.ErrInvalidJobStatus)
	}

	// Add log if any
	if log != "" {
		if err := tx.Log(context.Background(), j.id, log); err != nil {
			return err
		}
	}

	// Update job status
	//
	if err := tx.SetStatus(context.Background(), []jobq.ID{j.id}, jobq.JobStatusDone); err != nil {
		return err
	}

	return tx.Commit()
}

func (j *Job) Fail(log string) error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	// Job status must be reserved
	//
	status, err := tx.GetStatus(context.Background(), j.ID())
	if err != nil {
		return err
	}
	if status != jobq.JobStatusReserved {
		return fmt.Errorf("%w: job must be in reserved state to be relased", jobq.ErrInvalidJobStatus)
	}

	// Add log if any
	if log != "" {
		if err := tx.Log(context.Background(), j.id, log); err != nil {
			return err
		}
	}

	// Get job info
	//
	infos, err := tx.GetInfos(context.Background(), []jobq.ID{j.id})
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return jobq.ErrJobNotFound
	}
	info := infos[0]

	retries := info.Retries + 1
	// TODO: Update job Retries

	if retries > info.Options.MaxRetries {
		if err := tx.SetStatus(context.Background(), []jobq.ID{j.id}, jobq.JobStatusCanceled); err != nil {
			return err
		}

		return nil
	}

	// calculate backoff delay
	//
	delay := info.Options.MinBackOff
	for i := 0; i < int(retries); i++ {
		delay *= 2
		if delay > info.Options.MaxBackOff {
			delay = info.Options.MinBackOff
		}
	}

	status, err = j.service.pqRepo.Push(context.Background(), j.topic, info.Priority, j.id, time.Now().Add(delay))
	if err != nil {
		return err
	}

	// Update job status
	//
	if err := tx.SetStatus(context.Background(), []jobq.ID{j.id}, status); err != nil {
		return err
	}

	return tx.Commit()
}

func (j *Job) Cancel(log string) error {
	return fmt.Errorf("not yet implemented")
}
