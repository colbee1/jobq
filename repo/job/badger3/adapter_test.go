package badger3

import (
	"context"
	"testing"
	"time"

	"github.com/colbee1/jobq"
	"github.com/stretchr/testify/require"
)

var jobTestRepo *Adapter
var jobTestTopic = jobq.Topic("test")

func getRepo() (*Adapter, error) {
	if jobTestRepo == nil {
		dbPath := "../../../test/db/job-badger3-test"

		if r, err := New(dbPath, RepositoryOptions{DropDB: true}); err != nil {
			return nil, err
		} else {
			jobTestRepo = r
		}
	}

	return jobTestRepo, nil
}

func TestAdapterNew(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)
}

func TestAdapterCreate(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	opts := jobq.DefaultJobOptions
	opts.Name = "test job"
	opts.LogStatusChange = true
	jid, err := tx.Create(context.Background(), jobTestTopic, -100, opts, jobq.Payload("Hello world"))
	require.NoError(err)
	require.Equal(jobq.ID(1), jid)
	err = tx.Commit()
	require.NoError(err)
}

func TestAdapterRead(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	jobs, err := tx.Read(context.Background(), []jobq.ID{jobq.ID(1)})
	require.NoError(err)
	require.Len(jobs, 1)
	require.Equal(jobq.JobStatusCreated, jobs[0].Status)
	require.Equal(jobq.Weight(-100), jobs[0].Weight)
	require.Equal(jobTestTopic, jobs[0].Topic)
	require.Equal("test job", jobs[0].Options.Name)
}

func TestAdapterPayload(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	payload, err := tx.ReadPayload(context.Background(), jobq.ID(1))
	require.NoError(err)
	require.Equal("Hello world", string(payload))
}

func TestAdapterUpdate(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	err = tx.Update(context.Background(), []jobq.ID{jobq.ID(1)},
		func(job *jobq.JobInfo) error {
			job.Status = jobq.JobStatusReady
			return nil
		})
	require.NoError(err)
	err = tx.Commit()
	require.NoError(err)
}

func TestAdapterGetStatus2(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	jobs, err := tx.Read(context.Background(), []jobq.ID{jobq.ID(1)})
	require.NoError(err)
	require.Len(jobs, 1)
	require.Equal(jobq.JobStatusReady, jobs[0].Status)
}

func TestAdapterFindByStatus(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	createdJobs, err := tx.FindByStatus(context.Background(), jobq.JobStatusCreated, 0, 100)
	require.NoError(err)
	require.Len(createdJobs, 0)

	readyJobs, err := tx.FindByStatus(context.Background(), jobq.JobStatusReady, 0, 100)
	require.NoError(err)
	require.Len(readyJobs, 1)
}

func TestAdapterReadMultiple(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	infos, err := tx.Read(context.Background(), []jobq.ID{jobq.ID(1), jobq.ID(2)})
	require.NoError(err)
	require.Len(infos, 1)
	info := infos[0]
	expected := &jobq.JobInfo{
		ID:             jobq.ID(1),
		Topic:          jobq.Topic(jobTestTopic),
		Weight:         -100,
		Status:         jobq.JobStatusReady,
		DateCreated:    info.DateCreated,    // TODO:
		DateTerminated: info.DateTerminated, // TODO:
		DatesReserved:  info.DatesReserved,  // TODO:
		RetryCount:     0,
		Options: jobq.JobOptions{
			Name:            "test job",
			Timeout:         jobq.DefaultJobOptions.Timeout,
			DelayedAt:       time.Time{},
			MaxRetries:      jobq.DefaultJobOptions.MaxRetries,
			MinBackOff:      jobq.DefaultJobOptions.MinBackOff,
			MaxBackOff:      jobq.DefaultJobOptions.MaxBackOff,
			LogStatusChange: true,
		},
		Logs: info.Logs, // TODO:
	}
	require.Equal(expected, info)
}

func TestAdapterClose(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	err = repo.Close()
	require.NoError(err)
}
