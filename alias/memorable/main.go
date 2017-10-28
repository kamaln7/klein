package memorable

import (
	"math/rand"
	"strings"
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

// New initializes the alias generator and returns a new instance
func New(c *Config) *Provider {
	c.randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Provider{
		Config: c,
	}
}

// Generate returns a random alias
func (p *Provider) Generate() string {
	var (
		output = ""
		length = len(wordlist)
	)

	for i := 0; i < p.Config.Length; i++ {
		output += strings.Title(wordlist[p.Config.randSource.Intn(length)])
	}

	return output
}
