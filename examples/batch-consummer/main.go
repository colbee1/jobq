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
	jobRepo, err := job_adapter.New("../../test/db/examples/batch-consummer-job",
		job_adapter.RepositoryOptions{DropDB: true})
	if err != nil {
		return err
	}
	defer jobRepo.Close()

	topicRepo, err := topic_adapter.New("../../test/db/examples/batch-consummer-topic",
		topic_adapter.RepositoryOptions{DropDB: true, StatsCollector: true})
	if err != nil {
		return err
	}
	defer topicRepo.Close()

	jqSvc, err := jqs.New(jobRepo, topicRepo)
	if err != nil {
		return err
	}
	defer jqSvc.Close()

	// Start the batch consumer
	//
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go consumer(jqSvc, wg)

	// Simulate producer
	//
	producer(jqSvc)

	wg.Wait()

	// List job done
	//
	jids, err := jqSvc.FindByStatus(context.Background(), jobq.JobStatusDone, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Done\n", len(jids))

	// List job reserved
	//
	jids, err = jqSvc.FindByStatus(context.Background(), jobq.JobStatusReserved, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Reserved\n", len(jids))

	// List job canceled
	//
	jids, err = jqSvc.FindByStatus(context.Background(), jobq.JobStatusCanceled, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Canceled\n", len(jids))

	// List job delayed
	//
	jids, err = jqSvc.FindByStatus(context.Background(), jobq.JobStatusDelayed, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Delayed\n", len(jids))

	// List job Ready
	//
	jids, err = jqSvc.FindByStatus(context.Background(), jobq.JobStatusReady, 0, 1000)
	if err != nil {
		return err
	}
	fmt.Printf("\n%d jobs are in status Ready\n", len(jids))

	// Show topic stats
	//
	list, err := jqSvc.Topics(context.Background(), 0, 100)
	if err != nil {
		return err
	}
	for _, topic := range list {
		stats, err := jqSvc.TopicStats(context.Background(), topic)
		if err != nil {
			return err
		}
		fmt.Printf("\ntopics: %s, stats:%+v\n", topic, stats)
	}

	return nil
}

func producer(jq service.IJobService) {
	prng := rand.New(rand.NewSource(time.Now().UnixNano()))
	upper := math.MaxInt16
	lower := math.MinInt16

	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(prng.Intn(150)) * time.Millisecond)

		priority := jobq.Weight(prng.Intn(upper-lower) + lower) // -10 is a higher priority than 10
		jid, err := jq.Enqueue(context.Background(), Topic, priority, jobq.JobOptions{}, jobq.Payload("email@domain.tld"))
		if err != nil {
			fmt.Printf("Enqueue error: %v\n", err)
			continue
		}
		fmt.Printf("Job%d was pushed in topic %s with priority %d\n", jid, Topic, priority)
	}

}

func consumer(jq service.IJobService, wg *sync.WaitGroup) {
	defer wg.Done()

	batchTicker := time.NewTicker(2 * time.Second)
	batchSize := 25

	noMoreJobCounter := 0
	for noMoreJobCounter < 2 {
		t := <-batchTicker.C
		jobs, err := jq.Reserve(context.Background(), Topic, batchSize)
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
			info, _ := job.Unwrap()
			pri := info.Weight
			batch = append(batch, fmt.Sprintf("Job #%d(pri=%d)", job.ID(), pri))
			if err := job.Done(); err != nil {
				panic(err)
			}
		}

		fmt.Printf("\n%s: Process batch of %d jobs: %s\n\n", t.Format(time.RFC3339), len(batch), strings.Join(batch, ", "))
	}
}
