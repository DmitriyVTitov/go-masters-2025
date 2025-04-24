package concurrency

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	queue := NewQueue[int]()

	queue.Enqueue(10)
	queue.Enqueue(20)
	queue.Enqueue(30)

	val, ok := queue.Dequeue()
	if ok {
		t.Log("Dequeued:", val)
	}
}
