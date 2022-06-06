package migrations

import (
	"testing"
)

var emptyStore = ExampleStore{}

func TestStoreIsEmpty(t *testing.T) {
	empty := IsEmpty(emptyStore, table)

	if !empty {
		t.Fatalf(`wanted empty, but got %v`, empty)
	}
}

func TestStoreIsTracked(t *testing.T) {
	tracked := IsTracked(emptyStore, table)

	if tracked {
		t.Fatalf(`wanted tracked == false, but got %v`, tracked)
	}
}

func TestStoreVersion(t *testing.T) {
	version, tracked := Version(emptyStore, table)

	if version.Version != "" || tracked {
		t.Fatalf(`wanted version == "0" && tracked == false, but got (%v, %v)`, version.Version, tracked)
	}
}

func TestStoreIsUpToDate(t *testing.T) {
	upToDate := IsUpToDate(emptyStore, table, defaultList())

	if upToDate {
		t.Fatalf(`wanted upToDate == false, but got %v`, upToDate)
	}
}

func TestStartTrackingStore(t *testing.T) {
	state := StartTracking(emptyStore, table)

	if state {
		t.Fatalf(`wanted state == false, but got %v`, state)
	}
}

func TestStopTrackingStore(t *testing.T) {
	state := StopTracking(emptyStore, table)

	// Empty store is not tracked. So, StopTracking should succeed.
	if !state {
		t.Fatalf(`wanted state == false, but got %v`, state)
	}
}
