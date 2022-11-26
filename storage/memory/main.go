package memory

import (
	"github.com/kamaln7/klein/storage"
)

// Provider implements a temporary in-memory storage
type Provider struct {
	Config *Config

	urls map[string]string
}

// Config contains the configuration for the in-memory storage
type Config struct {
}

// ensure that the storage.Provider interface is implemented
var _ storage.Provider = new(Provider)

// New returns a new Provider instance
func New(c *Config) *Provider {
	return &Provider{
		Config: c,
		urls:   make(map[string]string),
	}
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	url, found := p.urls[alias]
	if !found {
		return "", storage.ErrNotFound
	}

	return url, nil
}

// Exists checks if there is a URL with the requested alias
func (p *Provider) Exists(alias string) (bool, error) {
	_, found := p.urls[alias]

	return found, nil
}

// Store creates a new short URL
func (p *Provider) Store(url, alias string, overwrite bool) error {
	_, found := p.urls[alias]
	if found && !overwrite {
		return storage.ErrAlreadyExists
	}

	p.urls[alias] = url
	return nil
}
