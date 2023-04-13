# jobq

!! Early development - work in progress - only for comments !!

jobq is a simple embeddable golang package to handle jobs/tasks/xxx asynchronously.

## Features:

- Job are enqueued in weighted (priority) queue named topic. Lower is the weight, higher is the priority.
- Job can be delayed before enqueued in topic.
- Job and Topic relies on repository adapter so they can be persisted.
- Topic creation is on the fly. Just push/pop on it.
- Simple topic stats when enabling stats collector.
- Job FSM is simple, only few status: Ready, Delayed, Reserved, Done or Canceled.
- Job API is simple, only few usual methods: Done(), Retry(), Cancel(), Payload(), Logf()
- Job retry defaults to user-configurable exponential backoff.

### In progress

- Auto restart job after application crash when using a durable job adapter.

### Planned

- Job timeout.
- Auto cleanup unused topics.

# Design

![design](design.svg)

## Available job repositories

- repo/job/memory - volatile in-memory storage. Does not support transaction.
- repo/job/badger3 - durable file system storage. Support transactions.

## Available topic repositories

- repo/pq/memory - volatile in-memory fast priority queues. Does not support transaction.
- repo/pq/badger3 - durable file system storage. Support transactions.

## Examples

See examples directory.
