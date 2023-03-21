package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

// Create creates and store a new job
func (adapter *Adapter) Create(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jo jobq.JobOptions) (jobq.JobID, error) {
	if topic == "" {
		return 0, jobq.ErrTopicIsInvalid
	}

	jid := jobq.JobID(adapter.jobSequence.Add(1))
	job := &jobq.Job{
		ID:          jid,
		Status:      jobq.JobStateReady,
		Priority:    pri,
		Topic:       topic,
		DateCreated: time.Now(),
		Options:     jo,
	}

	// repo model is same as domain model
	adapter.jobs[job.ID] = job

	return jid, nil
}
