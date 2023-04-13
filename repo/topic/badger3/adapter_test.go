package badger3

import (
	"context"
	"testing"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo/topic"
	"github.com/stretchr/testify/require"
)

var pqTestRepo *Adapter
var TopicTest = jobq.Topic("test")

func getRepo() (*Adapter, error) {
	if pqTestRepo == nil {
		dbPath := "../../../test/db/pq-badger3-test"

		if r, err := New(dbPath, topic.RepositoryOptions{DropDB: true, StatsCollector: true}); err != nil {
			return nil, err
		} else {
			pqTestRepo = r
		}
	}

	return pqTestRepo, nil
}

func TestAdapterNew(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)
}

func TestPushTopic(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	status, err := repo.Push(context.Background(), TopicTest, 0, jobq.ID(1), time.Time{})
	require.NoError(err)
	require.Equal(jobq.JobStatusReady, status)
}

func TestPopTopic(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	jids, err := repo.PopTopic(context.Background(), TopicTest, 10)
	require.NoError(err)
	require.NotNil(jids)
	require.Len(jids, 1)
	require.Equal(jobq.ID(1), jids[0])

	jids, err = repo.PopTopic(context.Background(), TopicTest, 10)
	require.NoError(err)
	require.NotNil(jids)
	require.Len(jids, 0)
}

func TestPushDelayed(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	status, err := repo.Push(context.Background(), TopicTest, 0, jobq.ID(1), time.Now().Add(2*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)

	status, err = repo.Push(context.Background(), TopicTest, 0, jobq.ID(2), time.Now().Add(3*time.Second))
	require.NoError(err)
	require.Equal(jobq.JobStatusDelayed, status)
}

func TestPopDelayed(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	jids, err := repo.PopDelayed(context.Background(), 10)
	require.NoError(err)
	require.NotNil(jids)
	require.Len(jids, 0)

	time.Sleep(2 * time.Second)

	jids, err = repo.PopDelayed(context.Background(), 10)
	require.NoError(err)
	require.NotNil(jids)
	require.Len(jids, 1)

	time.Sleep(1 * time.Second)

	jids, err = repo.PopDelayed(context.Background(), 10)
	require.NoError(err)
	require.NotNil(jids)
	require.Len(jids, 1)

	jids, err = repo.PopDelayed(context.Background(), 10)
	require.NoError(err)
	require.NotNil(jids)
	require.Len(jids, 0)

}

func TestAdapterTopics(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	ts, err := repo.Topics(context.Background(), 0, 100)
	require.NoError(err)
	require.Len(ts, 2) // [ "test", "delayed" ]
}

func TestAdapterStatsTopic(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	stats, err := repo.TopicStats(context.Background(), TopicTest)
	require.NoError(err)
	require.True(stats.PopDateFirst.After(stats.PushDateFirst))
	require.Equal(int64(1), stats.PushTotalCount)
	require.Equal(int64(1), stats.PopTotalCount)

	stats, err = repo.TopicStats(context.Background(), "delayed")
	require.NoError(err)
	require.True(stats.PopDateFirst.After(stats.PushDateFirst))
	require.Equal(int64(2), stats.PushTotalCount)
	require.Equal(int64(2), stats.PopTotalCount)
}

func TestAdapterClose(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	err = repo.Close()
	require.NoError(err)
}
