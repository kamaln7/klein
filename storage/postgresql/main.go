package postgresql

import (
	"crypto/tls"
	"errors"

	"github.com/jackc/pgx"
	pgxstdlib "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kamaln7/klein/alias"
	"github.com/kamaln7/klein/storage"
)

// Provider implements a redis-based storage system
type Provider struct {
	Config *Config

	db *sqlx.DB
}

// Config contains the configuration for the redis storage
type Config struct {
	Host, User, Password, Database, Table, SSLMode string

	Alias alias.Provider
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
	cc := pgx.ConnConfig{
		Host:     p.Config.Host,
		User:     p.Config.User,
		Password: p.Config.Password,
		Database: p.Config.Database,
	}

	// Copied from https://github.com/jackc/pgx/blob/f25025a5801f9c925f4a7ffea5636bf53755c67e/conn.go#L976-L997
	switch p.Config.SSLMode {
	case "disable":
		cc.UseFallbackTLS = false
		cc.TLSConfig = nil
		cc.FallbackTLSConfig = nil
	case "allow":
		cc.UseFallbackTLS = true
		cc.FallbackTLSConfig = &tls.Config{InsecureSkipVerify: true}
	case "prefer":
		cc.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		cc.UseFallbackTLS = true
		cc.FallbackTLSConfig = nil
	case "require":
		cc.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	case "verify-ca", "verify-full":
		cc.TLSConfig = &tls.Config{
			ServerName: cc.Host,
		}
	default:
		return errors.New("sslmode is invalid")
	}

	db := pgxstdlib.OpenDB(cc)
	p.db = sqlx.NewDb(db, "pgx")

	p.createTable()

	return nil
}

func (p *Provider) createTable() error {
	q := `
	CREATE TABLE IF NOT EXISTS {{ TableName }} (
		alias serial,
		url text,
		PRIMARY KEY( alias )
	);
	`
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	return "", nil
}

// Exists checks if there is a URL with the requested alias
func (p *Provider) Exists(alias string) (bool, error) {
	return false, nil
}

// Store creates a new short URL
func (p *Provider) Store(url, alias string) error {
	return nil
}
