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
	"github.com/kamaln7/klein/storage/file"
)

var (
	length       = flag.Int("length", 3, "code length")
	key          = flag.String("key", "", "upload API Key")
	root         = flag.String("root", "", "root redirect")
	path         = flag.String("path", "/srv/www/urls/", "path to urls")
	listenAddr   = flag.String("listenAddr", "127.0.0.1:5556", "listen address")
	publicURL    = flag.String("url", "http://127.0.0.1:5556/", "path to public facing url")
	notFoundPath = flag.String("template", "./404.html", "path to error template")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "[klein] ", log.Ldate|log.Ltime)

	notFoundHTML, err := ioutil.ReadFile(*notFoundPath)
	if err != nil {
		logger.Fatal(err)
		return
	}

	var authProvider auth.Provider
	if *key == "" {
		authProvider = unauthenticated.New()
	} else {
		authProvider = statickey.New(&statickey.Config{
			Key: *key,
		})
	}

	k := klein.New(&klein.Config{
		Alias: alphanumeric.New(&alphanumeric.Config{
			Length: *length,
		}),
		Auth: authProvider,
		Storage: file.New(&file.Config{
			Path: *path,
		}),
		Log: logger,

		ListenAddr:   *listenAddr,
		RootURL:      *root,
		PublicURL:    *publicURL,
		NotFoundHTML: notFoundHTML,
	})

	k.Serve()
}
