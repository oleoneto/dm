package migrations

import "fmt"

type MigrationList struct {
	head *Migration
	tail *Migration
	size int
}

func (List *MigrationList) Size() int {
	return List.size
}

func (List *MigrationList) IsEmpty() bool {
	return List.size == 0
}

// GetHead - Returns the node at the start of the list.
func (List *MigrationList) GetHead() *Migration {
	return List.head
}

// GetTail - Returns the node at the end of the list.
func (List *MigrationList) GetTail() *Migration {
	return List.tail
}

// Insert - Adds a new node to the end of the list.
func (List *MigrationList) Insert(node *Migration) {
	if List.head == nil {
		List.head = node
		List.tail = node
	} else {
		node.previous = List.tail
		List.tail.next = node
		List.tail = node
	}

	List.size += 1
}

// Remove - Excludes a node from the list.
func (List *MigrationList) Remove(identifier string) {
	curr := List.head

	if curr.Version == identifier || curr.Name == identifier {
		List.head = curr.next
		List.size -= 1
		return
	}

	for curr != nil {
		if curr.Version == identifier || curr.Name == identifier {

			if curr.previous != nil {
				curr.previous.next = nil
			}

			if curr.next != nil {
				curr.previous.next = curr.next
			}

			if curr == List.tail {
				List.tail = curr.previous
			}

			List.size -= 1
			return
		}

		curr = curr.next
	}
}

// Reverse - Traverses list and swaps the direction of the list.
func (List *MigrationList) Reverse() {
	var prev *Migration
	curr := List.head

	// No head. An empty list.
	if curr == nil {
		return
	}

	// Swap pointers
	for curr != nil {
		next := curr.next
		curr.next = prev
		prev = curr
		curr = next
	}

	List.tail = List.head
	List.head = prev
}

// Display - Traverses list and prints all of its elements.
func (List *MigrationList) Display() {
	curr := List.head

	for curr != nil {
		format := fmt.Sprintf("%+v", curr.Name)
		if curr.next != nil {
			format += " -> "
		}
		fmt.Print(format)

		curr = curr.next
	}

	fmt.Println()
}

// Find - Traverses the list in search for a given node.
func (List *MigrationList) Find(identifier string) (MigrationList, bool) {
	curr := List.head
	var sequence MigrationList

	if curr == nil {
		return sequence, false
	}

	for curr != nil {
		sequence.Insert(&Migration{
			Changes:  curr.Changes,
			Engine:   curr.Engine,
			FileName: curr.FileName,
			Id:       curr.Id,
			Name:     curr.Name,
			Schema:   curr.Schema,
			Version:  curr.Version,
		})

		if curr.Version == identifier || curr.Name == identifier {
			return sequence, true
		}

		curr = curr.next
	}

	return sequence, false
}

// FromMap - Inserts map elements into the list.
func (List *MigrationList) FromMap(m map[string]Migration) {
	for _, value := range m {
		List.Insert(&value)
	}
}

// ToMap - Transforms the list into a map.
func (List *MigrationList) ToMap() map[string]Migration {
	m := map[string]Migration{}

	curr := List.head

	for curr != nil {
		m[curr.Version] = *curr
		curr = curr.next
	}

	return m
}

// ToSlice - Transforms the list into a slice.
func (List *MigrationList) ToSlice() Migrations {
	m := Migrations{}

	curr := List.head

	for curr != nil {
		m = append(m, *curr)
		curr = curr.next
	}

	return m
}

func (List *MigrationList) Description() string {
	if List.size == 0 {
		return "No migrations."
	}

	return fmt.Sprintf("%v migrations", List.size)
}
