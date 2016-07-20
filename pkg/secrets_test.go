package pkg

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/registry-token-refresher/pkg/credentials"
	"github.com/deis/registry-token-refresher/pkg/k8s"
	"k8s.io/kubernetes/pkg/api"
	apierrors "k8s.io/kubernetes/pkg/api/errors"
)

var dockerConfig = credentials.DockerConfig{Hostname: "test", Token: "testtoken"}

func TestCreateSecretErr(t *testing.T) {
	expectedErr := errors.New("get secret error")
	secretsClient := &k8s.FakeSecrets{
		FnCreate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, expectedErr
		},
	}
	err := CreateSecret(secretsClient, dockerConfig)
	assert.Err(t, err, expectedErr)
	assert.Equal(t, len(secretsClient.CreateCalls), 1, "number of create secret calls")
}

func TestCreateSecretSuccess(t *testing.T) {
	secretsClient := &k8s.FakeSecrets{
		FnCreate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, nil
		},
	}
	err := CreateSecret(secretsClient, dockerConfig)
	assert.NoErr(t, err)
	assert.Equal(t, len(secretsClient.CreateCalls), 1, "number of create secret calls")
}

func TestCreateSecretAlreadyExists(t *testing.T) {
	alreadyExistErr := apierrors.NewAlreadyExists(api.Resource("tests"), "1")
	secretsClient := &k8s.FakeSecrets{
		FnCreate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, alreadyExistErr
		},
		FnUpdate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, nil
		},
	}
	err := CreateSecret(secretsClient, dockerConfig)
	assert.NoErr(t, err)
	assert.Equal(t, len(secretsClient.CreateCalls), 1, "number of create secret calls")
	assert.Equal(t, len(secretsClient.UpdateCalls), 1, "number of update secret calls")
}

func TestUpdateSecretErr(t *testing.T) {
	expectedErr := errors.New("get secret error")
	secretsClient := &k8s.FakeSecrets{
		FnUpdate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, expectedErr
		},
	}
	err := UpdateSecret(secretsClient, dockerConfig)
	assert.Err(t, err, expectedErr)
	assert.Equal(t, len(secretsClient.UpdateCalls), 1, "number of update secret calls")
}

func TestUpdateSecretSuccess(t *testing.T) {
	secretsClient := &k8s.FakeSecrets{
		FnUpdate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, nil
		},
	}
	err := UpdateSecret(secretsClient, dockerConfig)
	assert.NoErr(t, err)
	assert.Equal(t, len(secretsClient.UpdateCalls), 1, "number of update secret calls")
}

func TestUpdateSecretNotFound(t *testing.T) {
	notFoundErr := apierrors.NewNotFound(api.Resource("tests"), "1")
	secretsClient := &k8s.FakeSecrets{
		FnUpdate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, notFoundErr
		},
	}
	err := UpdateSecret(secretsClient, dockerConfig)
	assert.NoErr(t, err)
	assert.Equal(t, len(secretsClient.UpdateCalls), 1, "number of update secret calls")
}
