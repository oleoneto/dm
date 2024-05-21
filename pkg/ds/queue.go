package ds

import "fmt"

type Queue[T any] struct{ data []T }

// Size - Returns the number of elements in the queue.
func (Q *Queue[T]) Size() int { return len(Q.data) }

// IsEmpty - Returns true if the queue has no elements.
func (Q *Queue[T]) IsEmpty() bool { return Q.Size() == 0 }

// GetFront - Returns the item at the start of the queue.
func (Q *Queue[T]) GetFront() *T {
	if Q.IsEmpty() {
		return nil
	}

	return &Q.data[0]
}

// GetBack - Returns the item at the end of the queue.
func (Q *Queue[T]) GetBack() *T {
	if Q.IsEmpty() {
		return nil
	}
	return &Q.data[len(Q.data)-1]
}

// Dequeue - Returns the item at the start of the queue.
func (Q *Queue[T]) Dequeue() *T {
	if Q.IsEmpty() {
		return nil
	}

	n := Q.data[0]

	if len(Q.data) == 1 {
		n = Q.data[0]
		Q.data = []T{}
		return &n
	}

	var newData []T = Q.data[1:]

	Q.data = newData
	return &n
}

// Enqueue - Adds a new item to the end of the queue.
func (Q *Queue[T]) Enqueue(item T) { Q.data = append(Q.data, item) }

// Remove - Excludes a item from the queue.
func (Q *Queue[T]) Remove(matcher func(item T) bool) {
	if Q.IsEmpty() {
		return
	}

	for index, item := range Q.data {
		if matcher(item) {
			newData := append(Q.data[:index], Q.data[index+1:]...)
			Q.data = newData
			return
		}
	}
}

// Find - Walks the queue in search of a given item. Returns the given item.
func (Q *Queue[T]) Find(finderFunc func(item T) bool) *T {
	if Q.IsEmpty() {
		return nil
	}

	for _, item := range Q.data {
		if finderFunc(item) {
			return &item
		}
	}

	return nil
}

// FindSequence - Walks the queue in search of a given item. Returns a new sequence starting at the found item.
func (Q *Queue[T]) FindSequence(finderFunc func(item T) bool) *Queue[T] {
	if Q.IsEmpty() {
		return nil
	}

	var matchedIndex = -1
	for index, item := range Q.data {
		if finderFunc(item) {
			matchedIndex = index
			break
		}
	}

	if matchedIndex < 0 {
		return nil
	}

	var queue Queue[T]
	for i := matchedIndex; i < len(Q.data); i++ {
		queue.Enqueue(Q.data[i])
	}

	return &queue
}

func (Q *Queue[T]) String() string { return fmt.Sprintf("%v migrations", Q.Size()) }

func (Q *Queue[T]) RawData() *[]T { return &Q.data }

func (Q *Queue[T]) Reversed() Queue[T] {
	if Q.IsEmpty() {
		return Queue[T]{}
	}

	total := Q.Size()
	newData := make([]T, total)
	for index := total - 1; index >= 0; index-- {
		newData[total-index-1] = Q.data[index]
	}
	return *NewFromSlice(newData)
}

func (Q *Queue[T]) Reverse() { *Q = Q.Reversed() }

// MARK: Initializers

func NewQueue[T any]() *Queue[T] { return &Queue[T]{data: []T{}} }

func NewFromSlice[T any](data []T) *Queue[T] { return &Queue[T]{data: data} }
