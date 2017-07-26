package alphanumeric

import (
	"math/rand"
	"time"

	"github.com/kamaln7/klein/alias"
)

// Provider implements an alias generator
type Provider struct {
	Config *Config
}

// ensure that the storage.Provider interface is implemented
var _ alias.Provider = new(Provider)

// Config contains the configuration for the file storage
type Config struct {
	Length     int
	randSource *rand.Rand
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// New initializes the alias generator and returns a new instance
func New(c *Config) *Provider {
	c.randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Provider{
		Config: c,
	}
}

// Generate returns a random alias
func (p *Provider) Generate() string {
	b := make([]rune, p.Config.Length)
	for i := range b {
		b[i] = runes[p.Config.randSource.Intn(len(runes))]
	}

	return string(b)
}
