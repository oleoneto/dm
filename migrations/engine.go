package migrations

import (
	"io/fs"
	"regexp"
)

type Engine interface {
	FileLoader
	MigrationRunner
	Tracker
	Validator
}

type FileLoader interface {
	// LoadFiles - Loads files from the given directory matching the provided regex
	LoadFiles(string, *regexp.Regexp) []fs.FileInfo
}

type MigrationRunner interface {
	// Up - Runs migrations
	Up(changes MigrationList) error

	// Down - Reverts migrations
	Down(changes MigrationList) error
}

type Tracker interface {
	// ... -
	BuildMigrations([]fs.FileInfo) MigrationList

	// StartTracking - Prepares database for migration tracking
	StartTracking() error

	// StopTracking - Stops tracking database migrations
	StopTracking() error

	// Version - Return the version of the last applied migration. The returned boolean should indicate if the database is being tracked
	Version() (string, bool)

	// IsUpToDate - Indicator of whether migrations are current or up-to-date
	IsUpToDate(changes MigrationList) bool

	// IsTracked - Indicator of whether the database is being managed by this tool
	IsTracked() bool

	// IsEmpty - Indicator of whether the database has any migrations
	IsEmpty() bool

	// AppliedMigrations - Returns all applied/saved migrations
	AppliedMigrations() map[string]Migration

	// PendingMigrations - Returns all non-applied/saved migrations.
	PendingMigrations() map[string]Migration
}

type Validator interface {
	// Validate - Given a set of migrations, this method should return whether or not the migrations are valid.
	Validate(changes MigrationList) bool
}
