package jqs

import (
	"context"
	"testing"
	"time"

	repo_job "github.com/colbee1/jobq/repo/job/memory"
	repo_topic "github.com/colbee1/jobq/repo/topic/memory"
	"github.com/stretchr/testify/require"

	"github.com/colbee1/jobq"
)

func TestServiceWithMemoryRepos(t *testing.T) {
	require := require.New(t)

	jobRepo, err := repo_job.New()
	require.NoError(err)
	require.NotNil(jobRepo)
	defer jobRepo.Close()

	pqRepo, err := repo_topic.New()
	require.NoError(err)
	require.NotNil(pqRepo)
	defer pqRepo.Close()

	s, err := New(jobRepo, pqRepo)
	require.NoError(err)
	require.NotNil(s)

	ctx := context.Background()

	jid, err := s.Enqueue(ctx, "", -10, jobq.DefaultJobOptions, jobq.Payload{})
	require.NoError(err)
	require.Equal(jobq.ID(1), jid)

	reserved, err := s.Reserve(ctx, "", 10)
	require.NoError(err)
	require.Equal(1, len(reserved))
	job := reserved[0]
	err = job.Logf("message #1\n")
	require.NoError(err)

	reserved, err = s.Reserve(ctx, "", 10)
	require.NoError(err)
	require.Equal(0, len(reserved))

	err = job.Done()
	require.NoError(err)

	jid, err = s.Enqueue(ctx, "", 3, jobq.JobOptions{DelayedAt: time.Now().Add(3 * time.Second)}, jobq.Payload{})
	require.NoError(err)
	require.NotNil(jid)

	jid, err = s.Enqueue(ctx, "", 2, jobq.JobOptions{DelayedAt: time.Now().Add(2 * time.Second)}, jobq.Payload{})
	require.NoError(err)
	require.NotNil(jid)

	// Because job is delayed check there is nothing to reserve
	reserved, err = s.Reserve(ctx, "", 10)
	require.NoError(err)
	require.Equal(0, len(reserved))

	time.Sleep(4 * time.Second)

	reserved, err = s.Reserve(ctx, "", 10)
	require.NoError(err)
	require.Len(reserved, 2)
	for _, job := range reserved {
		job.Done()
	}
}
