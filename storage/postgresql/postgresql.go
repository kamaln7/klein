package postgresql

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	pgxstdlib "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kamaln7/klein/storage"
)

// Provider implements a PostgreSQL-based storage system
type Provider struct {
	Config *Config

	db *sqlx.DB
}

// Config contains the configuration for the PostgreSQL server and database
type Config struct {
	Host, User, Password, Database, Table, SSLMode string
	Port                                           int32
}

// ensure that the storage.Provider interface is implemented
var _ storage.Provider = new(Provider)

// database url type
type url struct {
	ID         int
	URL, Alias string
}

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

// Init sets up the PostgreSQL database connection and creates the table if needed
func (p *Provider) Init() error {
	cc := pgx.ConnConfig{
		Host:     p.Config.Host,
		Port:     uint16(p.Config.Port),
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

	// create table if it doesn't already exist
	if err := p.createTable(); err != nil {
		return err
	}

	return nil
}

func (p *Provider) fillInTableName(query string) string {
	return fmt.Sprintf(query, p.Config.Table)
}

func (p *Provider) createTable() error {
	q := p.fillInTableName(`
	create table if not exists %s (
		id serial,
		alias text unique not null,
		url text not null,
		primary key( id )
	)`)

	_, err := p.db.Exec(q)
	return err
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	u := &url{}

	q := p.fillInTableName("select * from %s where alias = $1")
	err := p.db.Get(u, q, alias)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.ErrNotFound
		}

		return "", err
	}

	return u.URL, nil
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
	q := p.fillInTableName("insert into %s (url, alias) values ($1, $2)")
	_, err := p.db.Exec(q, url, alias)

	if err, ok := err.(*pgx.PgError); ok && err.Code == "23505" {
		return storage.ErrAlreadyExists
	}

	return err
}

func (p *Provider) DeleteURL(alias string) error {
	q := p.fillInTableName("DELETE FROM %s WHERE alias = $1")
	_, err := p.db.Exec(q, alias)
	if err != nil {
		return err
	}
	return nil
}
