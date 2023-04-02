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

	jids, err := tx.FindByStatus(context.Background(), jobq.JobStatusCreated, 0, 0)
	require.NoError(err)
	require.Len(jids, 1)

	err = tx.Update(context.Background(), []jobq.ID{jobq.ID(jid)},
		func(job *jobq.JobInfo) error {
			job.Status = jobq.JobStatusReady
			return nil
		})
	require.NoError(err)

	jids, err = tx.FindByStatus(context.Background(), jobq.JobStatusCreated, 0, 0)
	require.NoError(err)
	require.Len(jids, 0)

	err = tx.Logf(context.Background(), jid, "log test")
	require.NoError(err)

	jobs, err := tx.Read(context.Background(), []jobq.ID{jid})
	require.NoError(err)
	require.Len(jobs, 1)
	job := jobs[0]

	require.Len(job.Logs, 2)
	require.Equal("test", string(job.Topic))
	require.Equal(-10, int(job.Priority))
	require.Equal(jobq.JobStatusReady, job.Status)
	require.Equal(jobq.JobStatusReady, job.Status)
}
