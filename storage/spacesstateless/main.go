package spacesstateless

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kamaln7/klein/storage"
	cache "github.com/patrickmn/go-cache"
)

// Provider implements an in memory storage that persists on DigitalOcean Spaces
type Provider struct {
	Config *Config

	spaces *s3.S3
	cache  *cache.Cache
}

// Config contains the configuration for the file storage
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Space     string
	Path      string

	CacheDuration time.Duration
}

// ensure that the storage.Provider interface is implemented
var _ storage.Provider = new(Provider)

// New returns a new Provider instance
func New(c *Config) (*Provider, error) {
	spacesSession := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Endpoint:    aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", c.Region)),
		Region:      aws.String("us-east-1"), // Needs to be us-east-1, or it'll fail.
	})
	spaces := s3.New(spacesSession)

	p := &Provider{
		Config: c,

		spaces: spaces,
	}

	if c.CacheDuration != 0 {
		p.cache = cache.New(c.CacheDuration, c.CacheDuration/2)
	}

	return p, nil
}

func (p *Provider) aliasFullPath(alias string) string {
	prefix := ""
	if p.Config.Path != "" {
		prefix = fmt.Sprintf("%s/", strings.TrimSuffix(p.Config.Path, "/"))
	}

	return fmt.Sprintf("%s%s", prefix, alias)
}

func (p *Provider) getFromSpaces(alias string) (string, error) {
	output, err := p.spaces.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.Config.Space),
		Key:    aws.String(p.aliasFullPath(alias)),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return "", storage.ErrNotFound
			case "InvalidAccessKeyId":
				log.Printf("storage/spaces-stateless: invalid access key, could not access spaces")
				return "", aerr
			default:
				return "", aerr
			}
		}

		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(output.Body)

	return buf.String(), nil
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	if p.cache == nil {
		return p.getFromSpaces(alias)
	}

	cachedURL, isCached := p.cache.Get(alias)
	if isCached {
		return cachedURL.(string), nil
	}

	url, err := p.getFromSpaces(alias)
	if err != nil {
		return "", err
	}

	p.cache.Set(alias, url, cache.DefaultExpiration)
	return url, nil
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
func (p *Provider) Store(url, alias string, overwrite string) error {
	exists, err := p.Exists(alias)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrAlreadyExists
	}

	object := s3.PutObjectInput{
		Body:   strings.NewReader(url),
		Bucket: aws.String(p.Config.Space),
		Key:    aws.String(p.aliasFullPath(alias)),
	}

	_, err = p.spaces.PutObject(&object)
	if err != nil {
		return err
	}

	if p.cache != nil {
		p.cache.Set(alias, url, cache.DefaultExpiration)
	}
	return nil
}
