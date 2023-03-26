package memory

import (
	"context"
	"testing"

	"github.com/colbee1/jobq"
	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	require := require.New(t)

	repo, err := New()
	require.NoError(err)
	require.NotNil(repo)
	defer repo.Close()

	require.False(repo.Durable())

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)

	_, err = tx.Create(context.Background(), "", 0, jobq.DefaultJobOptions, nil)
	require.ErrorIs(err, jobq.ErrTopicIsMissing)

	jOpts := jobq.DefaultJobOptions
	jOpts.LogStatusChange = true
	jid, err := tx.Create(context.Background(), "test", -10, jOpts, nil)
	require.NoError(err)
	require.Equal(jobq.ID(1), jid)

	jids, err := tx.ListByStatus(context.Background(), jobq.JobStatusCreated, 0, 0)
	require.NoError(err)
	require.Len(jids, 1)

	err = tx.SetStatus(context.Background(), []jobq.ID{jid}, jobq.JobStatusReady)
	require.NoError(err)

	jids, err = tx.ListByStatus(context.Background(), jobq.JobStatusCreated, 0, 0)
	require.NoError(err)
	require.Len(jids, 0)

	err = tx.Log(context.Background(), jid, "log test")
	require.NoError(err)

	logs, err := tx.Logs(context.Background(), jid)
	require.NoError(err)
	require.Len(logs, 3)

	infos, err := tx.GetInfos(context.Background(), []jobq.ID{jid})
	require.NoError(err)
	require.Len(infos, 1)

	info := infos[0]
	require.Equal("test", string(info.Topic))
	require.Equal(-10, int(info.Priority))
	require.Equal(jobq.JobStatusReady, info.Status)

	status, err := tx.GetStatus(context.Background(), jid)
	require.NoError(err)
	require.Equal(status, info.Status)
}
