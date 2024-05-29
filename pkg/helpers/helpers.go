package helpers

// FindLeftSequence - Walks the list in search of a given item. Returns a new sequence starting at the found item.
func FindLeftSequence[T any](data []T, finderFunc func(item T) bool) []T {
	if len(data) == 0 {
		return []T{}
	}

	var matchedIndex = -1
	for index, item := range data {
		if finderFunc(item) {
			matchedIndex = index
			break
		}
	}

	if matchedIndex < 0 {
		return []T{}
	}

	var list []T
	for i := 0; i < matchedIndex+1; i++ {
		list = append(list, data[i])
	}

	return list
}

// FindRightSequence - Walks the list in search of a given item. Returns a new sequence starting at the found item.
func FindRightSequence[T any](data []T, finderFunc func(item T) bool) []T {
	if len(data) == 0 {
		return []T{}
	}

	var matchedIndex = -1
	for index, item := range data {
		if finderFunc(item) {
			matchedIndex = index
			break
		}
	}

	if matchedIndex < 0 {
		return []T{}
	}

	var list []T
	for i := matchedIndex; i < len(data); i++ {
		list = append(list, data[i])
	}

	return list
}

// Applies a `transformer` function to every element in a list.
//
// Usage:
//
//	Map([]int{3, 4}, func(index, n int) int { return n * n }) // -> [9, 16]
func Map[A any, B any](collection []A, transformFunc func(int, A) B) []B {
	result := make([]B, len(collection))

	for index, item := range collection {
		result[index] = transformFunc(index, item)
	}

	return result
}
