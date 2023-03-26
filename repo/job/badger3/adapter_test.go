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

		if r, err := New(dbPath, Options{DropAll: true}); err != nil {
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
	tx.Commit()
}

func TestAdapterGetStatus(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	status, err := tx.GetStatus(context.Background(), jobq.ID(1))
	require.NoError(err)
	require.Equal(jobq.JobStatusCreated, status)
}

func TestAdapterSetStatus(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	err = tx.SetStatus(context.Background(), []jobq.ID{jobq.ID(1)}, jobq.JobStatusReady)
	require.NoError(err)
	tx.Commit()
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

	status, err := tx.GetStatus(context.Background(), jobq.ID(1))
	require.NoError(err)
	require.Equal(jobq.JobStatusReady, status)
}

func TestAdapterGetOptions(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	options, err := tx.GetOptions(context.Background(), jobq.ID(1))
	require.NoError(err)
	require.Equal("test job", options.Name)
	require.Equal(true, options.LogStatusChange)
}

func TestAdapterListByStatus(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	createdJobs, err := tx.ListByStatus(context.Background(), jobq.JobStatusCreated, 0, 100)
	require.NoError(err)
	require.Len(createdJobs, 0)

	readyJobs, err := tx.ListByStatus(context.Background(), jobq.JobStatusReady, 0, 100)
	require.NoError(err)
	require.Len(readyJobs, 1)
}

func TestAdapterGetPriority(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	pri, err := tx.GetPriority(context.Background(), jobq.ID(1))
	require.NoError(err)
	require.Equal(jobq.Priority(-100), pri)
}

func TestAdapterGetInfo(t *testing.T) {
	require := require.New(t)

	repo, err := getRepo()
	require.NoError(err)
	require.NotNil(repo)

	tx, err := repo.NewTransaction()
	require.NoError(err)
	require.NotNil(tx)
	defer tx.Close()

	infos, err := tx.GetInfos(context.Background(), []jobq.ID{jobq.ID(1), jobq.ID(2)})
	require.NoError(err)
	require.Len(infos, 1)
	info := infos[0]
	expected := &jobq.JobInfo{
		ID:             jobq.ID(1),
		Topic:          jobq.Topic(jobTestTopic),
		Priority:       -100,
		Status:         jobq.JobStatusReady,
		DateCreated:    info.DateCreated,    // TODO:
		DateTerminated: info.DateTerminated, // TODO:
		DateReserved:   info.DateReserved,   // TODO:
		Retries:        0,
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
