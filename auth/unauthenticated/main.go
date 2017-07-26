package unauthenticated

import (
	"net/http"

	"github.com/kamaln7/klein/auth"
)

// Provider implements an alias generator
type Provider struct{}

// ensure that the storage.Provider interface is implemented
var _ auth.Provider = new(Provider)

// New initializes the alias generator and returns a new instance
func New() *Provider {
	return &Provider{}
}

// Authenticate lets everyone go through
func (p *Provider) Authenticate(w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, nil
}
