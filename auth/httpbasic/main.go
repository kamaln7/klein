package httpbasic

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
	Username, Password string
}

// ensure that the auth.Provider interface is implemented
var _ auth.Provider = new(Provider)

// New initializes the auth provider and returns a new instance
func New(c *Config) *Provider {
	return &Provider{
		Config: c,
	}
}

// Authenticate makes sure the right credentials are passed
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (bool, error) {
	if p.Config.Username == "" || p.Config.Password == "" {
		return false, nil
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	username, password, authOK := r.BasicAuth()
	if authOK == false {
		return false, nil
	}

	if username != p.Config.Username || password != p.Config.Password {
		return false, nil
	}

	return true, nil
}
