package storage_test

import (
	"testing"

	. "github.com/kamaln7/klein/storage"
	"github.com/kamaln7/klein/storage/memory"
)

func TestProviders(t *testing.T) {
	providers := []Provider{
		newMemoryTestProvider(),
	}

	for _, p := range providers {
		testProvider(p, t)
	}
}

func testProvider(p Provider, t *testing.T) {
	t.Helper()
	var err error

	// test creating a new URL
	url := "http://example.com"
	alias := "example"
	err = p.Store(url, alias)
	if err != nil {
		t.Error("Couldn't store a new URL")
	}

	// test alias conflict
	err = p.Store(url, alias)
	if err != ErrAlreadyExists {
		t.Error("Couldn't handle storing a new URL with an existing alias properly")
	}
}

func newMemoryTestProvider() Provider {
	return memory.New(&memory.Config{})
}
