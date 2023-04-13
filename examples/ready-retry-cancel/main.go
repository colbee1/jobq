package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/job"
	job_adapter "github.com/colbee1/jobq/repo/job/badger3"
	topic_adapter "github.com/colbee1/jobq/repo/topic/badger3"
	"github.com/colbee1/jobq/service"
	"github.com/colbee1/jobq/service/jqs"
)

const Topic = jobq.Topic("demo")

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	jaOptions := job_adapter.RepositoryOptions{DropDB: true}
	jobRepo, err := job_adapter.New("../../test/db/examples/ready-retry-cancel-job", jaOptions)
	if err != nil {
		return err
	}
	defer jobRepo.Close()

	taOptions := topic_adapter.RepositoryOptions{
		StatsCollector: true,
		DropDB:         true,
	}
	topicRepo, err := topic_adapter.New("../../test/db/examples/ready-retry-cancel-topic", taOptions)
	if err != nil {
		return err
	}
	defer topicRepo.Close()

	jq, err := jqs.New(jobRepo, topicRepo)
	if err != nil {
		return err
	}
	defer jq.Close()

	// Start consumer
	//
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go consumer(jq, wg)

	// Enqueue a new job with delay
	//
	jid, err := jq.Enqueue(context.Background(), Topic, 0, jobq.JobOptions{DelayedAt: time.Now().Add(2 * time.Second)}, jobq.Payload("email@domain.tld"))
	if err != nil {
		fmt.Printf("Enqueue error: %v\n", err)
	}
	fmt.Printf("Job%d was pushed in topic %s\n", jid, Topic)

	showByStatus(jobRepo)

	wg.Wait()

	showByStatus(jobRepo)

	return nil
}

func consumer(jq service.IJobService, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Wait at most 5 seconds for 1 job...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	jobs, err := jq.Reserve(ctx, Topic, 1)
	if err != nil {
		fmt.Printf("Oups, reservation error: %v\n", err)
		return
	}
	if len(jobs) == 0 {
		fmt.Printf("Got zero job in 5 seconds")
		return
	}
	fmt.Printf("Got jobs: %v, retry it in 2 seconds\n", jobs)
	jobs[0].Retry(2 * time.Second)
	fmt.Printf("Job will be retried in 2 seconds...\n")

	fmt.Printf("Wait at most 5 seconds for 3 new jobs...\n")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	jobs, err = jq.Reserve(ctx, Topic, 3)
	if err != nil {
		fmt.Printf("Reserve error: %v\n", err)
		return
	}
	if len(jobs) == 0 {
		fmt.Printf("Got zero job in 5 seconds")
		return
	}
	fmt.Printf("Got job: %v, then cancel\n", jobs)
	job := jobs[0]
	job.Cancel()

	time.Sleep(500 * time.Microsecond)
	jq.Reset(context.Background(), []jobq.ID{job.ID()})

	fmt.Printf("Wait at most 2 seconds for 10 jobs...\n")
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	jobs, err = jq.Reserve(ctx, Topic, 10)
	if err != nil {
		panic(err)
	}
	if len(jobs) != 1 {
		panic("len(jobs) must be 1")
	}
}

func showByStatus(repo job.IJobRepository) error {
	tx, err := repo.NewTransaction()
	if err != nil {
		return err
	}
	defer tx.Close()

	for _, status := range []jobq.Status{jobq.JobStatusReady, jobq.JobStatusReserved, jobq.JobStatusDelayed, jobq.JobStatusDone, jobq.JobStatusCanceled} {
		if jids, err := tx.FindByStatus(context.Background(), status, 0, 100); err != nil {
			return err
		} else {
			fmt.Printf("%d jobs are in status %s: %v\n", len(jids), status.String(), jids)
		}
	}

	return nil
}
