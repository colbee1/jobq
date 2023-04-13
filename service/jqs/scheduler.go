package jqs

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/colbee1/jobq"
)

func (s *Service) scheduleDelayedJobs(running *atomic.Bool) {
	defer running.Store(false)

	for {

		jids, err := s.topicRepo.PopDelayed(context.Background(), 10) // TODO: Make it configurable
		if err != nil {
			fmt.Println(err)
			break
		}

		if len(jids) == 0 {
			break
		}

		tx, err := s.jobRepo.NewTransaction()
		if err != nil {
			fmt.Println(err)
			break
		}
		defer tx.Close()

		err = tx.Update(context.Background(), jids,
			func(job *jobq.JobInfo) error {
				status, err := s.topicRepo.Push(context.Background(), job.Topic, job.Weight, job.ID, time.Time{})
				if err != nil {
					return err
				}

				job.Status = status

				return nil
			})
		if err != nil {
			fmt.Println(err)
			break
		}
		if err := tx.Commit(); err != nil {
			fmt.Println(err)
			break
		}

	}
}

func (s *Service) scheduleExpiredReservations(running *atomic.Bool) {
	defer running.Store(false)

	fmt.Printf("------> %v: schedule expired reservation started\n", time.Now())
	// time.Sleep(4 * time.Second)
}

func (s *Service) scheduler(exit chan struct{}) {
	defer s.wg.Done()

	var retryTimeoutRunning atomic.Bool
	retryTimeoutTicker := time.NewTicker(time.Second)
	defer retryTimeoutTicker.Stop()

	var dequeueDelayedRunning atomic.Bool
	dequeueDelayedTicker := time.NewTicker(time.Second)
	defer dequeueDelayedTicker.Stop()

	for {
		select {

		case <-exit:
			retryTimeoutTicker.Stop()
			dequeueDelayedTicker.Stop()

			return

		case <-retryTimeoutTicker.C:
			if retryTimeoutRunning.CompareAndSwap(false, true) {
				go s.scheduleExpiredReservations(&retryTimeoutRunning)
			}

		case <-dequeueDelayedTicker.C:
			if dequeueDelayedRunning.CompareAndSwap(false, true) {
				go s.scheduleDelayedJobs(&dequeueDelayedRunning)
			}

		}
	}
}
