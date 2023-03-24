package memory

import (
	"context"
	"time"

	"github.com/colbee1/jobq"
)

// Create creates a new job and returns it's ID
func (t *Transaction) Create(ctx context.Context, topic jobq.JobTopic, pri jobq.JobPriority, jo jobq.JobOptions, payload jobq.JobPayload) (jobq.JobID, error) {
	a := t.a

	if topic == "" {
		return 0, jobq.ErrTopicIsMissing
	}

	jid := jobq.JobID(a.jobSequence.Add(1))
	model := &modelJob{
		id:       jid,
		topic:    topic,
		priority: pri,
		status:   jobq.JobStatusCreated,
		payload:  payload,
		options:  jo,
		info:     jobq.JobInfo{DateCreated: time.Now()},
		logs:     []string{},
	}

	if model.options.LogStatusChange {
		model.log("job created")
	}

	a.jobs[jid] = model
	t.needCommit = true

	return jid, nil
}
