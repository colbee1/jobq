package memory

type JobMinHeap []*JobItem

func (h JobMinHeap) Len() int           { return len(h) }
func (h JobMinHeap) Less(i, j int) bool { return h[i].heapPriority < h[j].heapPriority }
func (h JobMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *JobMinHeap) Push(x any) {
	*h = append(*h, x.(*JobItem))
}

func (h *JobMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]

	return x
}

func (h JobMinHeap) Peek() any {
	return h[0]
}
