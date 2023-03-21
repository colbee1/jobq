# jobq

jobq is a simple golang package to handle job/taks/... queue with priority and optional persistance

!! This a work in progress !!

## Example

```go
import (
    "github.com/colbee1/jobq"
    jqs "github.com/colbee1/jobq/service"
    jq_job_repo "github.com/colbee1/jobq/repo/job/memory"
    jq_pq_repo "github.com/colbee1/jobq/repo/pq/memory"
)

func main() {
    jobRepo := jq_job_repo.New()
    pqRepo := jq_pq_repo.New()
    jq, err := jqs.New(jobRepo, pqRepo)

    // Add a new job in topic "catalog-import"

    err := jq.Enqueue(context.Backgroun, "catalog-import", 0, jobq.JobOptions{
        Payload: []byte{"/path/to/file/data.csv"},
    })

    // Add a new job in topic "welcome" but delay it's execution in at least one hour

    err := jq.Enqueue(context.Backgroun, "welcome", 0, jobq.JobOptions{
        DelayedAt: time.Now().Add(time.Hour),
        Payload: []byte{"foo@bar.com"},
    })
}

func JobScheduler(jq jqs.JobServiceInterface) {
    // Reserve at most 10 waiting jobs on topic "catalog-import"
    jobs, err := jq.Reserve(context.Background, "catalog-import", 10)

    // Dispatch job to workers
    ...
}

func JobWorker(job *jqs.Job) {
    data := job.Payload()

    // do some work

    // you can add message(s) to job
    job.AppendMessage("hello world\n")

    // If all goes well
    job.Done()

    // or on error
    job.Error()

    // or reschedule the job in 30 minutes
    job.Release(time.Now().Add(30*time.Minute))

}
```
