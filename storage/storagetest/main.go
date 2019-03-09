package storagetest

import (
	"testing"

	"github.com/kamaln7/klein/storage"
)

// RunBasicTests run a basic test suite that should work on all storage providers
func RunBasicTests(p storage.Provider, t *testing.T) {
	var err error

	url := "http://example.com"
	alias := "example"

	t.Run("store new url", func(t *testing.T) {
		err = p.Store(url, alias)
		if err != nil {
			t.Error("couldn't store a new URL")
		}
	})

	t.Run("check existance of alias", func(t *testing.T) {
		exists, err := p.Exists(alias)
		if err != nil {
			t.Error(err)
		}

		if !exists {
			t.Error("expected alias exists got otherwise")
		}
	})

	t.Run("attempt to overwrite existing alias", func(t *testing.T) {
		err = p.Store(url, alias)
		if err != storage.ErrAlreadyExists {
			t.Error("couldn't handle storing a new URL with an existing alias properly")
		}
	})

	t.Run("look up existing alias", func(t *testing.T) {
		storedURL, err := p.Get(alias)
		if err != nil {
			t.Error("couldn't look up an existing alias")
		}
		if storedURL != url {
			t.Error("got a wrong url when looking up an alias")
		}
	})

	t.Run("look up nonexistant alias", func(t *testing.T) {
		// look up inexistent alias
		_, err = p.Get("1234567890")
		if err != storage.ErrNotFound {
			t.Error("didn't get the correct error looking up inexistent alias")
		}
	})
}
