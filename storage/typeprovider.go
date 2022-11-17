package storage

import (
	"errors"
)

// A Provider implements all the necessary functions for a storage backend for URLs
type Provider interface {
	Get(alias string) (string, error)
	Exists(alias string) (bool, error)
	Store(url, alias string, overwrite string) error
}

// Errors
var (
	ErrNotFound      = errors.New("URL does not exist")
	ErrAlreadyExists = errors.New("alias already exists")
)
