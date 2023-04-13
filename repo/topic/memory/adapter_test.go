package memory

import (
	"context"
	"testing"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/topic"

	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	require := require.New(t)

	repo, err := New(topic.RepositoryOptions{StatsCollector: true})
	require.NoError(err)
	require.NotNil(repo)
	defer repo.Close()

	ctx := context.Background()
	topic := jobq.Topic("test")

	status, err := repo.Push(ctx, topic, 0, jobq.ID(1), time.Time{})
	require.NoError(err)
	require.Equal(jobq.JobStatusReady, status)

	jids, err := repo.PopTopic(ctx, topic, 10)
	require.NoError(err)
	require.Len(jids, 1)
	require.Equal(jobq.ID(1), jids[0])

	// Add new delayed jobs.
	now := time.Now()
	status, err = repo.Push(ctx, topic, 0, jobq.ID(2), now.Add(4*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	status, err = repo.Push(ctx, topic, 0, jobq.ID(3), now.Add(2*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	status, err = repo.Push(ctx, topic, 0, jobq.ID(4), now.Add(3*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	// Poping right now on delayed queue should returns empty array
	jids, err = repo.PopDelayed(ctx, 10)
	require.NoError(err)
	require.Len(jids, 0)

	// Wait after delayed jobs expiration
	time.Sleep(5 * time.Second)

	// Poping now on delayed queue should returns all jobs in delayed order.
	jids, err = repo.PopDelayed(ctx, 10)
	require.NoError(err)
	require.Len(jids, 3)
	require.Equal([]jobq.ID{jobq.ID(3), jobq.ID(4), jobq.ID(2)}, jids)

	require.False(repo.Durable())

	time.Sleep(100 * time.Millisecond) // wait a moment because stats collecting is asynchronous
	stats, err := repo.TopicStats(context.Background(), topic)
	require.NoError(err)
	require.Equal(stats.PushTotalCount, stats.PopTotalCount)

	stats, err = repo.TopicStats(context.Background(), "delayed")
	require.NoError(err)
	require.Equal(stats.PushTotalCount, stats.PopTotalCount)
}
