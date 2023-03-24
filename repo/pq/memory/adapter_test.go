package memory

import (
	"context"
	"testing"
	"time"

	"github.com/colbee1/jobq"

	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	require := require.New(t)

	repo, err := New()
	require.NoError(err)
	require.NotNil(repo)
	defer repo.Close()

	ctx := context.Background()
	topic := jobq.JobTopic("test")

	status, err := repo.Push(ctx, topic, 0, jobq.JobID(1), time.Time{})
	require.NoError(err)
	require.Equal(jobq.JobStatusReady, status)

	l, err := repo.Len(ctx, topic)
	require.NoError(err)
	require.Equal(1, l)

	jids, err := repo.Pop(ctx, topic, 10)
	require.NoError(err)
	require.Len(jids, 1)
	require.Equal(jobq.JobID(1), jids[0])

	// Let priority reflect delayed time to pop in the same order and be able to check it
	now := time.Now()
	status, err = repo.Push(ctx, topic, 4, jobq.JobID(2), now.Add(4*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	status, err = repo.Push(ctx, topic, 2, jobq.JobID(3), now.Add(2*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	status, err = repo.Push(ctx, topic, 3, jobq.JobID(4), now.Add(3*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	// because jobs are delayed, queue length must be empty
	l, err = repo.Len(ctx, topic)
	require.NoError(err)
	require.Equal(0, l)

	// wait more than the max delayed job
	time.Sleep(5 * time.Second)

	// now jobs should be availables
	jids, err = repo.Pop(ctx, topic, 10)
	require.NoError(err)
	require.Len(jids, 3)
	require.Equal([]jobq.JobID{jobq.JobID(3), jobq.JobID(4), jobq.JobID(2)}, jids)

	require.False(repo.Durable())
}
