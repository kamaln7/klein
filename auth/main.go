package auth

import (
	"net/http"
)

// A Provider implements all the necessary functions for an authentication system
type Provider interface {
	Authenticate(w http.ResponseWriter, r *http.Request) (bool, error)
}
