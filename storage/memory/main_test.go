package memory

import (
	"testing"

	"github.com/kamaln7/klein/storage/storagetest"
)

func TestProvider(t *testing.T) {
	p := New(&Config{})

	storagetest.RunBasicTests(p, t)
}
