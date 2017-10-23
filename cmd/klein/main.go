package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/kamaln7/klein"
	"github.com/kamaln7/klein/alias/alphanumeric"
	"github.com/kamaln7/klein/auth"
	"github.com/kamaln7/klein/auth/statickey"
	"github.com/kamaln7/klein/auth/unauthenticated"
	"github.com/kamaln7/klein/storage"
	"github.com/kamaln7/klein/storage/bolt"
	"github.com/kamaln7/klein/storage/file"
)

var (
	length       = flag.Int("length", 3, "code length")
	key          = flag.String("key", "", "upload API Key")
	root         = flag.String("root", "", "root redirect")
	filepath     = flag.String("file.path", "", "path to urls")
	boltpath     = flag.String("bolt.path", "", "path to bolt db file")
	listenAddr   = flag.String("listenAddr", "127.0.0.1:5556", "listen address")
	publicURL    = flag.String("url", "http://127.0.0.1:5556/", "path to public facing url")
	notFoundPath = flag.String("template", "", "path to error template")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "[klein] ", log.Ldate|log.Ltime)

	if *filepath != "" && *boltpath != "" {
		logger.Fatalln("cannot use both file-based and boltdb-based storage")
	}

	notFoundHTML := []byte("404 not found")
	if *notFoundPath != "" {
		var err error
		notFoundHTML, err = ioutil.ReadFile(*notFoundPath)
		if err != nil {
			logger.Fatal(err)
			return
		}
	}

	var authProvider auth.Provider
	if *key == "" {
		authProvider = unauthenticated.New()
	} else {
		authProvider = statickey.New(&statickey.Config{
			Key: *key,
		})
	}

	var storage storage.Provider
	switch {
	case *filepath != "":
		storage = file.New(&file.Config{
			Path: *filepath,
		})
	case *boltpath != "":
		var err error
		storage, err = bolt.New(&bolt.Config{
			Path: *boltpath,
		})

		if err != nil {
			logger.Fatalf("could not open bolt database: %s\n", err.Error())
		}
	default:
		logger.Fatalln("please pass one storage engine")
	}

	k := klein.New(&klein.Config{
		Alias: alphanumeric.New(&alphanumeric.Config{
			Length: *length,
		}),
		Auth:    authProvider,
		Storage: storage,
		Log:     logger,

		ListenAddr:   *listenAddr,
		RootURL:      *root,
		PublicURL:    *publicURL,
		NotFoundHTML: notFoundHTML,
	})

	k.Serve()
}
