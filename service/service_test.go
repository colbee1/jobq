package service

import (
	"context"
	"testing"
	"time"

	jobq_job_repo "github.com/colbee1/jobq/repo/job/memory"
	jobq_pq_repo "github.com/colbee1/jobq/repo/pq/memory"
	"github.com/stretchr/testify/require"

	"github.com/colbee1/jobq"
)

func TestServiceWithMemoryRepos(t *testing.T) {
	require := require.New(t)

	jobRepo, err := jobq_job_repo.New()
	require.NoError(err)
	require.NotNil(jobRepo)
	defer jobRepo.Close()

	pqRepo, err := jobq_pq_repo.New()
	require.NoError(err)
	require.NotNil(pqRepo)
	defer pqRepo.Close()

	s, err := New(jobRepo, pqRepo)
	require.NoError(err)
	require.NotNil(s)

	ctx := context.Background()

	jid, err := s.Enqueue(ctx, "", -10, jobq.JobOptions{Name: "test"})
	require.NoError(err)
	require.Equal(jobq.JobID(1), jid)

	jobs, err := s.Find(ctx, []jobq.JobID{jid})
	require.NoError(err)
	require.Equal(1, len(jobs))
	require.Equal("test", jobs[0].Options.Name)
	require.Equal(jobq.JobStateReady, jobs[0].Status)

	reserved, err := s.Reserve(ctx, "", 10)
	require.NoError(err)
	require.Equal(1, len(reserved))
	job := reserved[0]
	err = job.AppendMessage(ctx, "message #1\n")
	require.NoError(err)

	reserved, err = s.Reserve(ctx, "", 10)
	require.NoError(err)
	require.Equal(0, len(reserved))

	jobs, err = s.Find(ctx, []jobq.JobID{jid})
	require.NoError(err)
	require.Equal(1, len(jobs))
	require.Equal("test", jobs[0].Options.Name)
	require.Equal(jobq.JobStateReserved, jobs[0].Status)
	require.Equal("message #1\n", jobs[0].Message)

	err = job.Done(ctx)
	require.NoError(err)

	jobs, err = s.FindByStatus(ctx, jobq.JobStateReady, 0, 0)
	require.NoError(err)
	require.Equal(0, len(jobs))

	jid, err = s.Enqueue(ctx, "", 3, jobq.JobOptions{DelayedAt: time.Now().Add(3 * time.Second)})
	require.NoError(err)
	require.NotNil(jid)

	jid, err = s.Enqueue(ctx, "", 2, jobq.JobOptions{DelayedAt: time.Now().Add(2 * time.Second)})
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
		job.Done(ctx)
	}
}
