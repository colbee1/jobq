package jqs

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/job"
	"github.com/colbee1/jobq/service"
)

type Job struct {
	service      *Service
	id           jobq.ID
	dateReserved time.Time
}

func newJob(service *Service, id jobq.ID, dateReserved time.Time) *Job {
	return &Job{
		service:      service,
		id:           id,
		dateReserved: dateReserved,
	}
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

	return nil, job.ErrJobNotFound
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

func (j *Job) cancel(tx job.IJobRepositoryTransaction) error {
	return tx.Update(context.Background(), []jobq.ID{jobq.ID(j.id)},
		func(job *jobq.JobInfo) error {
			job.Status = jobq.JobStatusCanceled
			job.DateTerminated = time.Now()

			return nil
		})
}

func (j *Job) Cancel() error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	if err := j.cancel(tx); err != nil {
		return err
	} else {
		return tx.Commit()
	}
}

func (j *Job) done(tx job.IJobRepositoryTransaction) error {
	return tx.Update(context.Background(), []jobq.ID{jobq.ID(j.id)},
		func(job *jobq.JobInfo) error {
			if job.Status != jobq.JobStatusReserved {
				return fmt.Errorf("%w: status must be Reserved to be switched to Done", service.ErrInvalidStatus)
			}

			if job.DatesReserved[len(job.DatesReserved)-1] != j.dateReserved {
				return fmt.Errorf("%w: Job has been requeued since it's reservation. Probably because of reservation timeout", service.ErrInvalidStatus)
			}

			job.Status = jobq.JobStatusDone
			job.DateTerminated = time.Now()

			return nil
		})
}

func (j *Job) Done() error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	if err := j.done(tx); err != nil {
		return err
	} else {
		return tx.Commit()
	}
}

func (j *Job) retry(tx job.IJobRepositoryTransaction, overrideBackoff time.Duration) error {
	return tx.Update(context.Background(), []jobq.ID{jobq.ID(j.id)},
		func(job *jobq.JobInfo) error {
			if job.Status != jobq.JobStatusReserved {
				return fmt.Errorf("%w: status must be Reserved to be retried", service.ErrInvalidStatus)
			}

			if job.RetryCount >= job.Options.MaxRetries {
				return j.cancel(tx)
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

			status, err := j.service.topicRepo.Push(context.Background(), job.Topic, job.Weight, job.ID, time.Now().Add(delay))
			if err != nil {
				return err
			}
			job.Status = status

			return nil
		})
}

func (j *Job) Retry(overrideBackoff time.Duration) error {
	tx, err := j.service.jobRepo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	if err := j.retry(tx, overrideBackoff); err != nil {
		return err
	} else {
		return tx.Commit()
	}
}
