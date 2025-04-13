package queue

import "fmt"

type node[T any] struct {
	Value T
	Next  *node[T]
}

// Element add in rear and remove from front
// Queue is a FIFO (First In First Out) data structure
type Queue[T any] struct {
	front *node[T] // head
	rear  *node[T] // back
	size  int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		front: nil,
		rear:  nil,
		size:  0,
	}
}

func (q *Queue[T]) IsEmpty() bool {
	return q.front == nil
}

func (q *Queue[T]) Enqueue(value T) {
	nn := &node[T]{Value: value, Next: nil}
	if q.front == nil {
		q.front = nn
	} else {
		q.rear.Next = nn
	}
	q.rear = nn
	q.size++
}

func (q *Queue[T]) Dequeue() error {
	if q.front == nil {
		return fmt.Errorf("queue is empty")
	}
	q.front = q.front.Next
	if q.front == nil {
		q.rear = nil
	}
	q.size--
	return nil
}

func (q *Queue[T]) Front() (T, error) {
	if q.front == nil {
		var zeroValue T
		return zeroValue, fmt.Errorf("queue is empty")
	}
	return q.front.Value, nil
}

func (q *Queue[T]) Size() int {
	return q.size
}
