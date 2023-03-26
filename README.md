# jobq

!! Early development - work in progress - only for comments !!

jobq is a simple embeddable golang package to handle jobs/tasks/xxx asynchronously.

## Features:

- Jobs are managed by fast priority queue.
- Jobs are pushed in named queues, called topics.
- Topic creation is on the fly. Just push on it.
- Job can be delayed in time.
- Job FSM is simple, only few status: Ready, Delayed, Reserved, Done or Cancel.
- Simple topic statistics.

### In progress

- Job that failed is automatically retried with definable backoff behavior. (IN PROGRESS)

### Planned

- Job can be persisted by using a "durable" repository backend.
  - Auto restart job after application crash.
- Job timeout.
- Topic with no more activity is automatically removed.

## Available job repository

- repo/job/memory (do not support transaction) - in memory job storage
- repo/job/badger3 (IN PROGRESS) - durable file system storage.

## Available priority queue repository

- repo/pq/memory - in memory priority queues

## Examples

See examples directory.
