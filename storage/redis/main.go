package redis

import (
	"github.com/kamaln7/klein/storage"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// Provider implements a redis-based storage system
type Provider struct {
	Config *Config
	pool   *pool.Pool
}

// Config contains the configuration for the redis storage
type Config struct {
	Address string
	Auth    string
	DB      int
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

// Init sets up the Redis pool
func (p *Provider) Init() error {
	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}

		if p.Config.Auth != "" {
			if err = client.Cmd("AUTH", p.Config.Auth).Err; err != nil {
				client.Close()
				return nil, err
			}
		}

		if err = client.Cmd("SELECT", p.Config.DB).Err; err != nil {
			client.Close()
			return nil, err
		}

		return client, nil
	}

	pool, err := pool.NewCustom("tcp", p.Config.Address, 10, df)
	if err != nil {
		return err
	}

	p.pool = pool
	return nil
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	r := p.pool.Cmd("GET", alias)
	if r.Err != nil {
		return "", r.Err
	}

	url, _ := r.Str()
	if url == "" {
		return "", storage.ErrNotFound
	}

	return url, nil
}

// Exists checks if there is a URL with the requested alias
func (p *Provider) Exists(alias string) (bool, error) {
	r, err := p.pool.Cmd("EXISTS", alias).Int()
	if err != nil {
		return false, err
	} else if r == 1 {
		return true, nil
	}

	return false, nil
}

// Store creates a new short URL
func (p *Provider) Store(url, alias string, overwrite bool) error {
	exists, err := p.Exists(alias)
	if err != nil {
		return err
	}
	if exists && !overwrite {
		return storage.ErrAlreadyExists
	}

	r := p.pool.Cmd("SET", alias, url)
	return r.Err
}
