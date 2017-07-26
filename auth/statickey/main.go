package statickey

import (
	"net/http"

	"github.com/kamaln7/klein/auth"
)

// Provider implements an alias generator
type Provider struct {
	Config *Config
}

// Config contains the config
type Config struct {
	Key string
}

// ensure that the storage.Provider interface is implemented
var _ auth.Provider = new(Provider)

// New initializes the alias generator and returns a new instance
func New(c *Config) *Provider {
	return &Provider{
		Config: c,
	}
}

// Authenticate makes sure the right key is passed
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (bool, error) {
	key := r.FormValue("key")

	if key == "" || key != p.Config.Key {
		return false, nil
	}

	return true, nil
}
