package jqs

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
	"github.com/colbee1/jobq/service"
)

type Job struct {
	service *Service
	id      jobq.ID
}

func (j *Job) ID() jobq.ID {
	return j.id
}

func (j *Job) Unwrap() (*jobq.JobInfo, error) {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	jobs, err := tx.Read(context.Background(), []jobq.ID{j.id})
	if err != nil {
		return nil, err
	}

	if len(jobs) > 0 {
		return jobs[0], nil
	}

	return nil, repo.ErrJobNotFound
}

func (j *Job) Payload() (jobq.Payload, error) {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Close()

	return tx.ReadPayload(context.Background(), j.id)
}

func (j *Job) Logf(format string, args ...any) error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	if err := tx.Logf(context.Background(), j.id, format, args...); err != nil {
		return err
	}

	return tx.Commit()
}

func (j *Job) Cancel() error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	err = tx.Update(context.Background(), []jobq.ID{jobq.ID(j.id)},
		func(job *jobq.JobInfo) error {
			job.Status = jobq.JobStatusCanceled
			job.DateTerminated = time.Now()

			return nil
		})
	if err != nil {
		return nil
	}

	return tx.Commit()
}

func (j *Job) Done() error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	err = tx.Update(context.Background(), []jobq.ID{jobq.ID(j.id)},
		func(job *jobq.JobInfo) error {
			if job.Status != jobq.JobStatusReserved {
				return fmt.Errorf("%w: status must be Reserved to be switched to Done", service.ErrInvalidStatus)
			}

			job.Status = jobq.JobStatusDone
			job.DateTerminated = time.Now()

			return nil
		})
	if err != nil {
		return nil
	}

	return tx.Commit()
}

func (j *Job) Retry(overrideBackoff time.Duration) error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	err = tx.Update(context.Background(), []jobq.ID{jobq.ID(j.id)},
		func(job *jobq.JobInfo) error {
			if job.Status != jobq.JobStatusReserved {
				return fmt.Errorf("%w: status must be Reserved to be retried", service.ErrInvalidStatus)
			}

			if job.RetryCount >= job.Options.MaxRetries {
				job.Status = jobq.JobStatusCanceled
				return nil
			} else {
				job.RetryCount++
			}

			// Calculate backoff delay
			//
			delay := job.Options.MinBackOff
			if overrideBackoff == 0 {
				for i := 0; i < int(job.RetryCount); i++ {
					delay *= 2
					if delay > job.Options.MaxBackOff {
						delay = job.Options.MinBackOff
					}
				}
			} else {
				delay = overrideBackoff
			}
			if delay < time.Second {
				delay = time.Second
			}

			status, err := j.service.pqRepo.Push(context.Background(), job.Topic, job.Priority, job.ID, time.Now().Add(delay))
			if err != nil {
				return err
			}
			job.Status = status

			return nil
		})
	if err != nil {
		return nil
	}

	return tx.Commit()
}
