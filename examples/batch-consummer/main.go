package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/colbee1/jobq"
	jobrepo "github.com/colbee1/jobq/repo/job/memory"
	pgrepo "github.com/colbee1/jobq/repo/pq/memory"
	"github.com/colbee1/jobq/services/jqs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	jobRepo, err := jobrepo.New()
	if err != nil {
		return err
	}
	defer jobRepo.Close()

	pqRepo, err := pgrepo.New()
	if err != nil {
		return err
	}
	defer pqRepo.Close()

	service, err := jqs.New(jobRepo, pqRepo)
	if err != nil {
		return err
	}
	defer service.Close()

	// Start the batch consumer
	//
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go consumer(service, wg)

	// Simulate producer
	//
	producer(service)

	wg.Wait()

	// List job done
	//
	jids, err := service.ListByStatus(context.Background(), jobq.JobStatusDone, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Done\n", len(jids))

	// List job reserved
	//
	jids, err = service.ListByStatus(context.Background(), jobq.JobStatusReserved, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Reserved\n", len(jids))

	return nil
}

func producer(jq jobq.IJobQueueService) {
	prng := rand.New(rand.NewSource(time.Now().UnixNano()))
	upper := math.MaxInt16
	lower := math.MinInt16

	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(prng.Intn(150)) * time.Millisecond)
		topic := jobq.JobTopic("sendmail")

		priority := jobq.JobPriority(prng.Intn(upper-lower) + lower) // -10 is a higher priority than 10
		jid, err := jq.Enqueue(context.Background(), topic, priority, jobq.JobOptions{}, jobq.JobPayload("email@domain.tld"))
		if err != nil {
			fmt.Printf("Enqueue error: %v\n", err)
			continue
		}
		fmt.Printf("Job%d was pushed in queue %s with priority %d\n", jid, topic, priority)
	}

}

func consumer(jq jobq.IJobQueueService, wg *sync.WaitGroup) {
	defer wg.Done()

	topic := jobq.JobTopic("sendmail")
	batchTicker := time.NewTicker(2 * time.Second)
	batchSize := 20

	noMoreJobCounter := 0
	for noMoreJobCounter < 2 {
		t := <-batchTicker.C
		jobs, err := jq.Reserve(context.Background(), topic, batchSize)
		if err != nil {
			fmt.Printf("Reserve error: %v\n", err)
			continue
		}

		if len(jobs) == 0 {
			noMoreJobCounter++
			continue
		}

		batch := []string{}
		for _, job := range jobs {
			pri, _ := job.Priority()
			batch = append(batch, fmt.Sprintf("Job%d(pri=%d)", job.ID(), pri))
			job.Done("")
		}

		fmt.Printf("\n%s: Process batch of %d jobs: %s\n\n", t.Format(time.RFC3339), len(batch), strings.Join(batch, ", "))
	}
}
