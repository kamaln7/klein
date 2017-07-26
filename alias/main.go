package alias

// A Provider implements all the necessary functions for an alias generator
type Provider interface {
	Generate() string
}
