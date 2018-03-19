package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/kamaln7/klein/alias"
	"github.com/kamaln7/klein/auth"
	"github.com/kamaln7/klein/storage"
)

// Klein is a URL shortener
type Klein struct {
	Config *Config
	mux    *http.ServeMux
}

// Config contains the necessary configuration to run the URL shortener
type Config struct {
	Alias        alias.Provider
	Auth         auth.Provider
	Storage      storage.Provider
	Log          *log.Logger
	NotFoundHTML []byte

	ListenAddr, PublicURL, RootURL string
}

// New returns a new Klein instance
func New(c *Config) *Klein {
	c.PublicURL = strings.TrimRight(c.PublicURL, "/") + "/"

	return &Klein{
		Config: c,
	}
}

// Serve starts Klein's HTTP server
func (b *Klein) Serve() {
	b.mux = http.NewServeMux()
	b.mux.HandleFunc("/", b.httpHandler)

	b.Config.Log.Printf("listening on %s\n", b.Config.ListenAddr)
	if err := http.ListenAndServe(b.Config.ListenAddr, b.mux); err != nil {
		b.Config.Log.Fatal(err)
	}
}

func (b *Klein) httpHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// root redirect & upload handlers
	if path == "/" {
		switch r.Method {
		case "GET":
			if b.Config.RootURL != "" {
				http.Redirect(w, r, b.Config.RootURL, 302)
			} else {
				b.notFound(w, r)
			}
		case "POST":
			b.create(w, r)
		}

		return
	}

	b.redirect(w, r, path[1:])
}

func (b *Klein) redirect(w http.ResponseWriter, r *http.Request, alias string) {
	url, err := b.Config.Storage.Get(alias)

	switch err {
	case nil:
	case storage.ErrNotFound:
		b.notFound(w, r)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}

	http.Redirect(w, r, url, 302)
}

func (b *Klein) create(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		url = r.FormValue("url")
	)

	// authenticate
	authed, err := b.Config.Auth.Authenticate(w, r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}
	if !authed {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("unauthenticated"))
		return
	}

	// validate input
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("you need to pass a url"))
		return
	}

	// set an alias
	alias := r.FormValue("alias")
	if alias == "" {
		exists := true
		for exists {
			alias = b.Config.Alias.Generate()
			exists, err = b.Config.Storage.Exists(alias)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error"))
				return
			}
		}
	} else {
		exists, err := b.Config.Storage.Exists(alias)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error"))
			return
		}

		if exists {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("code already exists"))
			return
		}
	}

	// store the URL
	err = b.Config.Storage.Store(url, alias)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(b.Config.PublicURL + alias))
}

func (b *Klein) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(b.Config.NotFoundHTML)
}
