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
		t.Error("couldn't store a new URL")
	}

	// test alias conflict
	err = p.Store(url, alias)
	if err != ErrAlreadyExists {
		t.Error("couldn't handle storing a new URL with an existing alias properly")
	}

	// look up alias
	storedUrl, err := p.Get(alias)
	if err != nil {
		t.Error("couldn't look up an existing alias")
	}
	if storedUrl != url {
		t.Error("got a wrong url when looking up an alias")
	}

	// look up inexistent alias
	_, err = p.Get("1234567890")
	if err != ErrNotFound {
		t.Error("couldn't look up an inexistent alias")
	}
}

func newMemoryTestProvider() Provider {
	return memory.New(&memory.Config{})
}
