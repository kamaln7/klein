package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/kamaln7/klein/storage"
)

// Provider implements a file-based storage system
type Provider struct {
	Config *Config
	mutex  sync.RWMutex
}

// Config contains the configuration for the file storage
type Config struct {
	Path string
}

// ensure that the storage.Provider interface is implemented
var _ storage.Provider = new(Provider)

// New returns a new Provider instance
func New(c *Config) *Provider {
	return &Provider{
		Config: c,
	}
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	alias = path.Base(alias)

	p.mutex.RLock()
	url, err := ioutil.ReadFile(filepath.Join(p.Config.Path, alias))
	p.mutex.RUnlock()
	if err != nil {
		return "", storage.ErrNotFound
	}

	return string(bytes.TrimSpace(url)), nil
}

// Exists checks if there is a URL with the requested alias
func (p *Provider) Exists(alias string) (bool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	_, err := os.Stat(filepath.Join(p.Config.Path, path.Base(alias)))
	return !os.IsNotExist(err), nil
}

// Store creates a new short URL
func (p *Provider) Store(url, alias string) error {
	exists, _ := p.Exists(alias)
	if exists {
		return storage.ErrAlreadyExists
	}

	p.mutex.Lock()
	err := ioutil.WriteFile(filepath.Join(p.Config.Path, alias), bytes.TrimSpace([]byte(url)), 0644)
	p.mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}
