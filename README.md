# jobq

!! Early development - work in progress - only for comments !!

jobq is an simple embeddable golang package to handle jobs/tasks/xxx asynchronously.

Expected features:

- Job is handled by a priority queue. (OK)
- Job is pushed in a named queue, called topic. (OK)
- Topic creation is automatic, just push on it. (OK)
- Job can be delayed. (OK)
- Job FSM is simple, only few status: Ready, Delayed, Reserved, Done or Cancel. (OK)
- Job that failed is automatically retried with definable backoff behavior. (IN PROGRESS)
- Job can be persisted by using a "durable" repository backend. (TODO)
- Auto restart job after application crash. (TODO)
- Job timeout. (TODO)

## Available job repository

- repo/job/memory (do not support transaction)
- repo/pq/badger3 (IN PROGRESS)

## Available priority queue repository

- repo/pq/memory (OK)

## Examples

See examples directory.
