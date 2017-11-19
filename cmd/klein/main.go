package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kamaln7/klein/alias"

	"github.com/kamaln7/klein"
	"github.com/kamaln7/klein/alias/alphanumeric"
	"github.com/kamaln7/klein/alias/memorable"
	"github.com/kamaln7/klein/auth"
	"github.com/kamaln7/klein/auth/httpbasic"
	"github.com/kamaln7/klein/auth/statickey"
	"github.com/kamaln7/klein/auth/unauthenticated"
	"github.com/kamaln7/klein/storage"
	"github.com/kamaln7/klein/storage/bolt"
	"github.com/kamaln7/klein/storage/file"
	"github.com/kamaln7/klein/storage/redis"
)

var (
	alphanumericLength = flag.Int("alphanumeric.length", -1, "alphanumeric code length")
	alphanumericAlpha  = flag.Bool("alphanumeric.alpha", true, "use letters in code")
	alphanumericNum    = flag.Bool("alphanumeric.num", true, "use numbers in code")
	memorablelength    = flag.Int("memorable.length", -1, "memorable word count")
	authkey            = flag.String("auth.key", "", "upload API Key")
	authusername       = flag.String("auth.username", "", "username for HTTP basic auth")
	authpassword       = flag.String("auth.password", "", "password for HTTP basic auth")
	root               = flag.String("root", "", "root redirect")
	filepath           = flag.String("file.path", "", "path to urls")
	boltpath           = flag.String("bolt.path", "", "path to bolt db file")
	redisaddress       = flag.String("redis.address", "", "address:port of redis instance")
	redisauth          = flag.String("redis.auth", "", "password to access redis")
	redisdb            = flag.Int("redis.db", 0, "db to select within redis")
	listenAddr         = flag.String("listenAddr", "127.0.0.1:5556", "listen address")
	publicURL          = flag.String("url", "", "path to public facing url")
	notFoundPath       = flag.String("template", "", "path to error template")
)

func exclusiveFlag(cb func(v ...interface{}), collisionErr, missingErr string, defaultVal interface{}, required bool, flags ...interface{}) {
	found := false
	for _, flag := range flags {
		if flag == defaultVal {
			continue
		}

		if found == true {
			cb(collisionErr)
		}

		found = true
	}

	if found == false && required {
		cb(missingErr)
	}
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "[klein] ", log.Ldate|log.Ltime)

	// exclusive flags
	exclusiveFlag(
		logger.Fatalln,
		"cannot use both file-based and boltdb-based storage",
		"please pass one storage provider",
		"",
		true,
		*filepath, *boltpath, *redisaddress,
	)
	exclusiveFlag(
		logger.Fatalln,
		"cannot use both alphanumeric and memorable alias providers",
		"please pass one alias provider",
		-1,
		true,
		*alphanumericLength, *memorablelength,
	)
	exclusiveFlag(
		logger.Fatalln,
		"cannot use both static key based auth and http basic auth",
		"please pass one auth provider",
		"",
		false,
		*authkey, *authusername,
	)

	// 404

	notFoundHTML := []byte("404 not found")
	if *notFoundPath != "" {
		var err error
		notFoundHTML, err = ioutil.ReadFile(*notFoundPath)
		if err != nil {
			logger.Fatal(err)
			return
		}
	}

	// auth

	var authProvider auth.Provider
	if *authkey != "" {
		authProvider = statickey.New(&statickey.Config{
			Key: *authkey,
		})
	} else if *authusername != "" {
		authProvider = httpbasic.New(&httpbasic.Config{
			Username: *authusername,
			Password: *authpassword,
		})
	} else {
		authProvider = unauthenticated.New()
	}

	// storage

	var storageProvider storage.Provider
	switch {
	case *filepath != "":
		storageProvider = file.New(&file.Config{
			Path: *filepath,
		})
	case *boltpath != "":
		var err error
		storageProvider, err = bolt.New(&bolt.Config{
			Path: *boltpath,
		})

		if err != nil {
			logger.Fatalf("could not open bolt database: %s\n", err.Error())
		}
	case *redisaddress != "":
		var err error
		storageProvider, err = redis.New(&redis.Config{
			Address: *redisaddress,
			Auth:    *redisauth,
			DB:      *redisdb,
		})

		if err != nil {
			logger.Fatalf("could not open redis database: %s\n", err.Error())
		}
	}

	// alias

	var aliasProvider alias.Provider
	switch {
	case *alphanumericLength != -1:
		var err error
		aliasProvider, err = alphanumeric.New(&alphanumeric.Config{
			Length: *alphanumericLength,
			Alpha:  *alphanumericAlpha,
			Num:    *alphanumericNum,
		})

		if err != nil {
			logger.Fatalf("could not select alphanumeric alias: %s\n", err.Error())
		}
	case *memorablelength != -1:
		aliasProvider = memorable.New(&memorable.Config{
			Length: *memorablelength,
		})
	}

	// url

	if *publicURL == "" {
		*publicURL = fmt.Sprintf("http://%s/", *listenAddr)
	}

	// klein

	k := klein.New(&klein.Config{
		Alias:   aliasProvider,
		Auth:    authProvider,
		Storage: storageProvider,
		Log:     logger,

		ListenAddr:   *listenAddr,
		RootURL:      *root,
		PublicURL:    *publicURL,
		NotFoundHTML: notFoundHTML,
	})

	k.Serve()
}
