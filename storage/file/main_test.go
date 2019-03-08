package file

import (
	"io/ioutil"
	"testing"

	"github.com/kamaln7/klein/storage/storagetest"
)

func TestProvider(t *testing.T) {
	dir, err := ioutil.TempDir("", "klein")
	if err != nil {
		t.Errorf("couldn't create temporary test dir: %v\n", err)
	}

	p := New(&Config{
		Path: dir,
	})

	storagetest.RunBasicTests(p, t)
}
