package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJobHeap1(t *testing.T) {
	require := require.New(t)

	h := newJobQueue()

	ji1 := &JobItem{heapPriority: 10, JobID: 1}
	h.Push(ji1)

	ji2 := &JobItem{heapPriority: 0, JobID: 2}
	h.Push(ji2)

	ji3 := &JobItem{heapPriority: -100, JobID: 3}
	h.Push(ji3)

	ji4 := &JobItem{heapPriority: -3, JobID: 4}
	h.Push(ji4)

	jis, err := h.Pop(10)
	require.NoError(err)
	require.Len(jis, 4)
	require.Equal([]*JobItem{ji3, ji4, ji2, ji1}, jis)
}

func TestJobHeap2(t *testing.T) {
	require := require.New(t)

	h := newJobQueue()

	ji1 := &JobItem{heapPriority: 1679394534, JobID: 1}
	h.Push(ji1)

	ji2 := &JobItem{heapPriority: 1679394530, JobID: 2}
	h.Push(ji2)

	ji3 := &JobItem{heapPriority: 1679394532, JobID: 3}
	h.Push(ji3)

	jis, err := h.Pop(10)
	require.NoError(err)
	require.Len(jis, 3)
	require.Equal([]*JobItem{ji2, ji3, ji1}, jis)
}
