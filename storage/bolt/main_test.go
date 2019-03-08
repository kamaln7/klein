package bolt

import (
	"io/ioutil"
	"testing"

	"github.com/kamaln7/klein/storage/storagetest"
)

func TestProvider(t *testing.T) {
	file, err := ioutil.TempFile("", "klein")
	if err != nil {
		t.Errorf("couldn't create temporary test file: %v\n", err)
	}

	p, err := New(&Config{
		Path: file.Name(),
	})
	if err != nil {
		t.Errorf("couldn't init bolt driver: %v\n", err)
	}

	storagetest.RunBasicTests(p, t)
}
