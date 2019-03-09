package bolt

import (
	"bytes"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kamaln7/klein/storage"
)

// Provider implements a file-based storage system
type Provider struct {
	Config *Config
	db     *bolt.DB
}

// Config contains the configuration for the file storage
type Config struct {
	Path string
}

// ensure that the storage.Provider interface is implemented
var _ storage.Provider = new(Provider)

// New returns a new Provider instance
func New(c *Config) (*Provider, error) {
	provider := &Provider{
		Config: c,
	}
	err := provider.Init()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Init sets up the BoltDB database
func (p *Provider) Init() error {
	db, err := bolt.Open(p.Config.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("klein"))

		return err
	})
	if err != nil {
		return err
	}

	p.db = db
	return nil
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	var url []byte

	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("klein"))
		url = bytes.TrimSpace(b.Get([]byte(alias)))
		if url == nil {
			return storage.ErrNotFound
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return string(url), nil
}

// Exists checks if there is a URL with the requested alias
func (p *Provider) Exists(alias string) (bool, error) {
	_, err := p.Get(alias)

	if err == storage.ErrNotFound {
		return false, nil
	}

	return true, err
}

// Store creates a new short URL
func (p *Provider) Store(url, alias string) error {
	exists, err := p.Exists(alias)
	if err != nil {
		return err
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	err = p.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("klein"))
		err := b.Put([]byte(alias), bytes.TrimSpace([]byte(url)))

		return err
	})

	return err
}
