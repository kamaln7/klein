package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/kamaln7/klein/storage/storagetest"
)

func TestProvider(t *testing.T) {
	redisPassword := "secret-password"

	redisServer, err := miniredis.Run()
	if err != nil {
		t.Errorf("couldn't start redis client: %v\n", err)
	}
	redisServer.RequireAuth(redisPassword)
	defer redisServer.Close()

	p, err := New(&Config{
		Address: redisServer.Addr(),
		DB:      5,
		Auth:    redisPassword,
	})
	if err != nil {
		t.Errorf("couldn't connect to redis server: %v\n", err)
	}

	storagetest.RunBasicTests(p, t)
}
