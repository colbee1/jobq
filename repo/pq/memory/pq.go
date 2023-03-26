package memory

import (
	"container/heap"
	"sync"
	"time"

	"github.com/colbee1/jobq"
	"github.com/colbee1/jobq/repo"
)

type (
	JobQueue struct {
		mu             sync.RWMutex
		queue          JobMinHeap
		dateCreated    time.Time
		dateLastPush   time.Time
		pushTotalCount int64
		maxQueueLen    int64
	}

	JobItem struct {
		heapPriority int64 // Used to order items (ie: timestamp)
		Topic        jobq.Topic
		Priority     jobq.Priority
		JobID        jobq.ID
	}
)

func newJobQueue() *JobQueue {
	jq := &JobQueue{
		queue:       JobMinHeap{},
		dateCreated: time.Now(),
	}
	heap.Init(&jq.queue)

	return jq
}

func (pq *JobQueue) Stats() repo.TopicStats {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return repo.TopicStats{
		DateCreated:          pq.dateCreated,
		DateLastPush:         pq.dateLastPush,
		PushTotalCount:       pq.pushTotalCount,
		MaxQueueLen:          pq.maxQueueLen,
		CurrentQueueCapacity: int64(cap(pq.queue)),
	}
}

func (pq *JobQueue) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return pq.queue.Len()
}

func (pq *JobQueue) Push(jitem *JobItem) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	heap.Push(&pq.queue, jitem)

	pq.pushTotalCount++
	pq.dateLastPush = time.Now()
	if l := int64(pq.queue.Len()); l > pq.maxQueueLen {
		pq.maxQueueLen = l
	}
}

func (pq *JobQueue) Pop(limit uint) ([]*JobItem, error) {
	resp := make([]*JobItem, 0, limit)
	if limit == 0 {
		return resp, nil
	}

	pq.mu.Lock()
	defer pq.mu.Unlock()

	for ; limit > 0 && pq.queue.Len() > 0; limit-- {
		item := heap.Pop(&pq.queue).(*JobItem)
		resp = append(resp, item)
	}

	return resp, nil
}

func (pq *JobQueue) Peek() *JobItem {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return pq.queue.Peek().(*JobItem)
}

// Recap purposes is avoiding continuously growing array
// by (temporary) reducing it's capacity.
func (pq *JobQueue) Recap(minToRecap int, spareFactor float32) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if spareFactor < 1 {
		spareFactor = 1.1
	}

	l := len(pq.queue)
	max := int(spareFactor * float32(l))
	cap := cap(pq.queue)
	if cap > minToRecap && cap > max {
		pq.queue = pq.queue[0:l:max]
	}
}
