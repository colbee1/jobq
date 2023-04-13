package memory

import (
	"container/heap"
	"sync"
	"time"

	"github.com/colbee1/jobq"
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
		Priority     jobq.Weight
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

func (pq *JobQueue) Pop() *JobItem {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.queue.Len() == 0 {
		return nil
	}

	return heap.Pop(&pq.queue).(*JobItem)
}

func (pq *JobQueue) Peek() *JobItem {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if v, ok := pq.queue.Peek().(*JobItem); ok {
		return v
	}

	return nil
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
