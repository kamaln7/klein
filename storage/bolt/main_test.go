package bolt

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kamaln7/klein/storage/storagetest"
)

func TestProvider(t *testing.T) {
	file, err := ioutil.TempFile("", "klein")
	if err != nil {
		t.Errorf("couldn't create temporary test file: %v\n", err)
	}
	defer os.Remove(file.Name())

	p, err := New(&Config{
		Path: file.Name(),
	})
	if err != nil {
		t.Errorf("couldn't init bolt driver: %v\n", err)
	}

	storagetest.RunBasicTests(p, t)
}
