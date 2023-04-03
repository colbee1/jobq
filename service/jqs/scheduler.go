package jqs

import (
	"context"
	"fmt"
	"time"

	"github.com/colbee1/jobq"
)

func (s *Service) scheduleDelayedJobs(exit chan struct{}) {
	defer s.wg.Done()

	secTicker := time.NewTicker(time.Second)
	defer secTicker.Stop()

loop:
	for {
		select {

		case <-exit:
			break loop

		case <-secTicker.C:
			for {

				jids, err := s.pqRepo.PopDelayed(context.Background(), 10) // TODO: Make it configurable
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
						status, err := s.pqRepo.Push(context.Background(), job.Topic, job.Priority, job.ID, time.Time{})
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
	}
}
