package memory

func NewTestProvider() *Provider {
	return &Provider{
		Config: &Config{},
	}
}
