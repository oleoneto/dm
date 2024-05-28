package runner_test

import (
	"io/fs"
	"regexp"
	"time"

	"github.com/oleoneto/dm/pkg/fsystem"
)

type PgOptions struct {
	Username       string
	Password       string
	Name           string
	MaxConnections string
	Port           uint32
}

type MockDirEntry struct {
	name       string
	isDir      bool
	mode       fs.FileMode
	modifiedAt time.Time
	sys        any
	size       int64
}

var _ fs.DirEntry = (*MockDirEntry)(nil)

func (e *MockDirEntry) Name() string { return e.name }

func (e MockDirEntry) IsDir() bool { return e.isDir }

func (e MockDirEntry) Type() fs.FileMode { return e.mode }

func (e *MockDirEntry) Info() (fs.FileInfo, error) { return e, nil }

func (f *MockDirEntry) Size() int64 { return f.size }

func (f *MockDirEntry) Mode() fs.FileMode { return f.mode }

func (f *MockDirEntry) ModTime() time.Time { return f.modifiedAt }

func (f *MockDirEntry) Sys() any { return f.sys }

// MARK: - File Loader

type Loader struct{ data []fs.DirEntry }

var _ fsystem.FileLoaderProtocol = (*Loader)(nil)

func (l Loader) LoadFiles(string, *regexp.Regexp) []fs.DirEntry { return l.data }
