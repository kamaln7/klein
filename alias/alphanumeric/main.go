package alphanumeric

import (
	"errors"
	"math/rand"
	"time"

	"github.com/kamaln7/klein/alias"
)

// Provider implements an alias generator
type Provider struct {
	Config *Config
	runes  []rune
}

// ensure that the storage.Provider interface is implemented
var _ alias.Provider = new(Provider)

// Config contains the configuration for the file storage
type Config struct {
	Length     int
	Alpha      bool
	Num        bool
	randSource *rand.Rand
}

var alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var num = []rune("0123456789")

// New initializes the alias generator and returns a new instance
func New(c *Config) (*Provider, error) {
	c.randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

	provider := &Provider{
		Config: c,
	}
	err := provider.Init()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Init sets up the alphanumeric alias
func (p *Provider) Init() error {
	var runes []rune
	switch {
	case p.Config.Alpha == true:
		runes = append(runes, alpha...)
		fallthrough
	case p.Config.Num == true:
		runes = append(runes, num...)
	default:
		return errors.New("please specify at least alpha or numeric!")
	}

	p.runes = runes
	return nil
}

// Generate returns a random alias
func (p *Provider) Generate() string {
	b := make([]rune, p.Config.Length)
	for i := range b {
		b[i] = p.runes[p.Config.randSource.Intn(len(p.runes))]
	}

	return string(b)
}
