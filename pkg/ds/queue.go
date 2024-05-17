package ds

import "fmt"

type Queue[T any] struct{ data []T }

// Size - Returns the number of elements in the queue.
func (Q *Queue[T]) Size() int { return len(Q.data) }

// IsEmpty - Returns true if the queue has no elements.
func (Q *Queue[T]) IsEmpty() bool { return Q.Size() == 0 }

// GetFront - Returns the node at the start of the queue.
func (Q *Queue[T]) GetFront() *T {
	if Q.IsEmpty() {
		return nil
	}

	return &Q.data[0]
}

// GetBack - Returns the node at the end of the queue.
func (Q *Queue[T]) GetBack() *T {
	if Q.IsEmpty() {
		return nil
	}
	return &Q.data[len(Q.data)-1]
}

// Dequeue - Returns the node at the start of the queue.
func (Q *Queue[T]) Dequeue() *T {
	if Q.IsEmpty() {
		return nil
	}

	n := Q.data[0]

	if Q.Size() == 1 {
		Q.data = []T{}
		return &n
	}

	if Q.Size() == 2 {
		Q.data = []T{Q.data[1]}
		return &n
	}

	newData := append(Q.data[:1], Q.data[len(Q.data)-1:]...)
	Q.data = newData

	return &n
}

// Enqueue - Adds a new node to the end of the queue.
func (Q *Queue[T]) Enqueue(node T) { Q.data = append(Q.data, node) }

// Remove - Excludes a node from the queue.
func (Q *Queue[T]) Remove(matcher func(node T) bool) {
	for index, node := range Q.data {
		if matcher(node) {
			newData := append(Q.data[:index], Q.data[index+1:]...)
			Q.data = newData
			return
		}
	}
}

// Find - Walks the queue in search of a given node.
func (Q *Queue[T]) Find(finder func(node T) bool) (queue Queue[T], isFound bool) {
	var matchedIndex = -1
	for index, node := range Q.data {
		if finder(node) {
			matchedIndex = index
			break
		}
	}

	if matchedIndex < 0 {
		return
	}

	for i := matchedIndex; i < len(Q.data); i++ {
		queue.Enqueue(Q.data[i])
	}

	return queue, true
}

func (Q *Queue[T]) Description() string {
	if Q.IsEmpty() {
		return "No migrations in list"
	}
	return fmt.Sprintf("%v migrations", Q.Size())
}

func (Q *Queue[T]) RawData() []T { return Q.data }

func (Q *Queue[T]) SetData(data []T) { Q.data = data }
