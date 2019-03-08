package storage_test

import (
	"testing"

	. "github.com/kamaln7/klein/storage"
	"github.com/kamaln7/klein/storage/memory"
)

// Separate methods so that we can run just one of them. Nice for debugging, VS Code gives us run test annotations:
// https://dmar.by/L8Du4Dr7
func TestMemoryProvider(t *testing.T) {
	testProvider(newMemoryTestProvider(), t)
}

func testProvider(p Provider, t *testing.T) {
	var err error

	url := "http://example.com"
	alias := "example"

	t.Run("store new url", func(t *testing.T) {
		err = p.Store(url, alias)
		if err != nil {
			t.Error("couldn't store a new URL")
		}
	})

	t.Run("attempt to overwrite existing alias", func(t *testing.T) {
		err = p.Store(url, alias)
		if err != ErrAlreadyExists {
			t.Error("couldn't handle storing a new URL with an existing alias properly")
		}
	})

	t.Run("look up existing alias", func(t *testing.T) {
		storedUrl, err := p.Get(alias)
		if err != nil {
			t.Error("couldn't look up an existing alias")
		}
		if storedUrl != url {
			t.Error("got a wrong url when looking up an alias")
		}
	})

	t.Run("look up nonexistant alias", func(t *testing.T) {
		// look up inexistent alias
		_, err = p.Get("1234567890")
		if err != ErrNotFound {
			t.Error("didn't get the correct error looking up inexistent alias")
		}
	})
}

func newMemoryTestProvider() Provider {
	return memory.New(&memory.Config{})
}
