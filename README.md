# jobq

!! Early development - work in progress - only for comments !!

jobq is a simple embeddable golang package to handle jobs/tasks/xxx asynchronously.

## Features:

- Job can be delayed in time.
- Jobs are managed by fast priority queue.
- Jobs can be persisted using durable repository.
- Jobs are pushed in named queues, called topics.
- Topic creation is on the fly. Just push/pop on it.
- Job FSM is simple, only few status: Ready, Delayed, Reserved, Done or Canceled.
- Simple topic statistics.
- Retry failed job with configurable exponential backoff.

### In progress

- Auto restart job after application crash. Only available with durable job repository.

### Planned

- Job timeout.
- Create LRU cache service in front of job repository to amortize read request.
- Auto cleanup unused topics.

## Available job repositories

- repo/job/memory - volatile in memory job storage. Does not support transaction.
- repo/job/badger3 - durable file system storage. Support transactions.

## Available priority queue repositories

- repo/pq/memory - in memory fast priority queues.

## Examples

See examples directory.
