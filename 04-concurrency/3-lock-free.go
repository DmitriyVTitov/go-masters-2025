package concurrency

import (
	"sync/atomic"
	"unsafe"
)

// ***
// Пример lock-free алгоритма на Go с использованием
// атомарных операций из пакета sync/atomic.
// ***

// Узел очереди
type Node[T any] struct {
	value T
	next  unsafe.Pointer
}

// Lock-free очередь FIFO
type LockFreeQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

// NewQueue создает новую очередь
func NewQueue[T any]() *LockFreeQueue[T] {
	dummy := &Node[T]{}
	queue := &LockFreeQueue[T]{head: unsafe.Pointer(dummy), tail: unsafe.Pointer(dummy)}
	return queue
}

// Enqueue добавляет элемент в очередь (без блокировки)
func (q *LockFreeQueue[T]) Enqueue(value T) {
	newNode := &Node[T]{value: value}

	for {
		tail := (*Node[T])(atomic.LoadPointer(&q.tail))
		next := (*Node[T])(atomic.LoadPointer(&tail.next))

		if tail == (*Node[T])(atomic.LoadPointer(&q.tail)) { // Проверяем, не изменился ли хвост
			if next == nil { // Если хвост указывает на последний элемент, добавляем новый узел
				if atomic.CompareAndSwapPointer(&tail.next, nil, unsafe.Pointer(newNode)) {
					atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(newNode))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			}
		}
	}
}

// Dequeue удаляет элемент из очереди (без блокировки)
func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	for {
		head := (*Node[T])(atomic.LoadPointer(&q.head))
		tail := (*Node[T])(atomic.LoadPointer(&q.tail))
		next := (*Node[T])(atomic.LoadPointer(&head.next))

		if head == (*Node[T])(atomic.LoadPointer(&q.head)) { // Проверяем, не изменился ли head
			if head == tail { // Если очередь пуста
				var zeroValue T
				return zeroValue, false
			}
			if next != nil {
				if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
					return next.value, true
				}
			}
		}
	}
}
