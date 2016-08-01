package credentials

import (
	"errors"
	"time"
)

type Parameters map[string]interface{}

// DockerConfig will hold information required for authentication to registry.
type DockerConfig struct {
	Token     string
	ExpiresAt time.Time
	Hostname  string
}

var errInvalidRegistry = errors.New("invalid registry type")

// DockerCredProvider defines methods that a docker credentials provider
// which uses short lived tokens for authentication must implement.
type DockerCredProvider interface {
	GetDockerCredentials() (DockerConfig, error)
	GetRefreshTime() time.Duration
}

// GetDockerCredentialsProvider returns an implementation of DockerCredProvider based
// on the registryLocation
func GetDockerCredentialsProvider(registryLocation string, params Parameters) (DockerCredProvider, error) {
	if registryLocation == "ecr" {
		return &ecrDockerCredProvider{params: params}, nil
	} else if registryLocation == "gcr" {
		return &gcrDockerCredProvider{params: params}, nil
	}
	return nil, errInvalidRegistry
}

// FakeDockerCredProvider is a mock function that can be swapped in for
// DockerCredProvider, so you can unit test your code.
type FakeDockerCredProvider struct {
	FnGetCreds   func() (DockerConfig, error)
	FnGetRefTime func() time.Duration
}

// GetDockerCredentials is the interface definition.
func (f *FakeDockerCredProvider) GetDockerCredentials() (DockerConfig, error) {
	return f.FnGetCreds()
}

// GetRefreshTime is the interface definition.
func (f *FakeDockerCredProvider) GetRefreshTime() time.Duration {
	return f.FnGetRefTime()
}
