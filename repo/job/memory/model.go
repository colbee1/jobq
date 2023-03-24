package memory

import (
	"time"

	"github.com/colbee1/jobq"
)

type modelJob struct {
	id       jobq.JobID
	topic    jobq.JobTopic
	priority jobq.JobPriority
	status   jobq.JobStatus
	payload  jobq.JobPayload
	options  jobq.JobOptions
	info     jobq.JobInfo
	logs     []string
}

func (mj *modelJob) log(msg string) {
	now := time.Now().Format(time.RFC3339Nano)
	mj.logs = append(mj.logs, now+": "+msg)
}
