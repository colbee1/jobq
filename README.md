# jobq

!! Early development - work in progress - only for comments !!

jobq is a simple embeddable golang package to handle jobs/tasks/xxx asynchronously.

## Features:

- Jobs are managed by fast priority queue.
- Jobs are pushed in named queues, called topics.
- Topic creation is on the fly. Just push on it.
- Job can be delayed in time.
- Job FSM is simple, only few status: Ready, Delayed, Reserved, Done or Canceled.
- Simple topic statistics.
- Jobs can be persisted using durable (and transactional) repository backend.

### In progress

- Retry failed job with configurable exponential backoff. (IN PROGRESS)

### Planned

- Auto restart job after application crash when using a durable job repository.
- Job timeout.
- Auto cleanup topic with no more activity.

## Available job repository

- repo/job/memory - volatile in memory job storage. Does not support transaction
- repo/job/badger3 - durable file system storage. In progresse 85% done.

## Available priority queue repository

- repo/pq/memory - in memory fast priority queues.

## Examples

See examples directory.
