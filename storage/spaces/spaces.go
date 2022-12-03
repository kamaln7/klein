package spaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kamaln7/klein/storage"
)

// Provider implements an in memory storage that persists on DigitalOcean Spaces
type Provider struct {
	Config *Config
	Spaces *s3.S3
	URLs   map[string]string
	mutex  sync.RWMutex
}

// Config contains the configuration for the file storage
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Space     string
	Path      string
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

	object := s3.GetObjectInput{
		Bucket: aws.String(c.Space),
		Key:    aws.String(c.Path),
	}

	urls := make(map[string]string)

	output, err := spaces.GetObject(&object)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return &Provider{
				Spaces: spaces,
				Config: c,
				URLs:   urls,
			}, nil
		}

		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(output.Body)

	err = json.Unmarshal(buf.Bytes(), &urls)
	if err != nil {
		return nil, err
	}

	return &Provider{
		Spaces: spaces,
		Config: c,
		URLs:   urls,
	}, nil
}

// Get attempts to find a URL by its alias and returns its original URL
func (p *Provider) Get(alias string) (string, error) {
	if url, exists := p.URLs[alias]; exists {
		return url, nil
	}

	return "", storage.ErrNotFound
}

// Exists checks if there is a URL with the requested alias
func (p *Provider) Exists(alias string) (bool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	_, exists := p.URLs[alias]
	return exists, nil
}

// Store creates a new short URL
func (p *Provider) Store(url, alias string) error {
	exists, _ := p.Exists(alias)
	if exists {
		return storage.ErrAlreadyExists
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.URLs[alias] = url

	body, err := json.Marshal(p.URLs)
	if err != nil {
		return err
	}

	object := s3.PutObjectInput{
		Body:   bytes.NewReader(body),
		Bucket: aws.String(p.Config.Space),
		Key:    aws.String(p.Config.Path),
	}
	_, err = p.Spaces.PutObject(&object)

	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) DeleteURL(alias string) error {
	exists, _ := p.Exists(alias)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if exists {
		delete(p.URLs, alias)
		body, err := json.Marshal(p.URLs)
		if err != nil {
			return err
		}
		object := s3.PutObjectInput{
			Body:   bytes.NewReader(body),
			Bucket: aws.String(p.Config.Space),
			Key:    aws.String(p.Config.Path),
		}
		_, err = p.Spaces.PutObject(&object)

		if err != nil {
			return err
		}
	}
	return nil
}
