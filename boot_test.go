package main

import (
	"errors"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/registry-token-refresher/pkg/credentials"
	"github.com/deis/registry-token-refresher/pkg/k8s"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func TestGetDiff(t *testing.T) {
	nsList := []api.Namespace{
		api.Namespace{ObjectMeta: api.ObjectMeta{Name: "app1"}},
		api.Namespace{ObjectMeta: api.ObjectMeta{Name: "app2"}},
	}
	added := getDiff(nsList)
	assert.Equal(t, len(added), 2, "number of namespaces added")
	assert.Equal(t, added, []string{"app1", "app2"}, "namespaces")
	nsList = append(nsList, api.Namespace{ObjectMeta: api.ObjectMeta{Name: "app3"}})
	added = getDiff(nsList)
	assert.Equal(t, len(added), 1, "number of namespaces added")
	assert.Equal(t, added, []string{"app3"}, "namespaces")
	added = getDiff(nsList)
	assert.Equal(t, len(added), 0, "number of namespaces added")
}

func TestTokenRefresherCredsErr(t *testing.T) {
	expectedErr := errors.New("get secret error")
	credProvider := &credentials.FakeDockerCredProvider{
		FnGetCreds: func() (credentials.DockerConfig, error) {
			return credentials.DockerConfig{}, expectedErr
		},
	}
	kubeClient := &k8s.FakeSecretsNamespacer{}
	appAddedCh := make(chan string)
	tokenRefErrCh := make(chan error)
	doneCh := make(chan struct{})
	go tokenRefresher(kubeClient, credProvider, appAddedCh, tokenRefErrCh, doneCh)
	timeoutCh := time.After(time.Second * 5)
	select {
	case err := <-tokenRefErrCh:
		assert.Err(t, err, expectedErr)
	case <-timeoutCh:
		t.Fatalf("expected error but didn't throw error")
	}
}

func TestTokenRefresherSecretErr(t *testing.T) {
	expectedErr := errors.New("get secret error")
	credProvider := &credentials.FakeDockerCredProvider{
		FnGetCreds: func() (credentials.DockerConfig, error) {
			return credentials.DockerConfig{}, nil
		},
		FnGetRefTime: func() time.Duration {
			return time.Minute * 1
		},
	}

	secretsClient := &k8s.FakeSecrets{
		FnCreate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, expectedErr
		},
	}
	kubeClient := &k8s.FakeSecretsNamespacer{
		Fn: func(string) kcl.SecretsInterface {
			return secretsClient
		},
	}
	appAddedCh := make(chan string)
	tokenRefErrCh := make(chan error)
	doneCh := make(chan struct{})
	go tokenRefresher(kubeClient, credProvider, appAddedCh, tokenRefErrCh, doneCh)
	appAddedCh <- "testapp"
	timeoutCh := time.After(time.Second * 5)
	select {
	case err := <-tokenRefErrCh:
		assert.Err(t, err, expectedErr)
	case <-timeoutCh:
		t.Fatalf("expected error but didn't throw error")
	}
	assert.Equal(t, len(secretsClient.CreateCalls), 1, "number of create secret calls")
}

func TestTokenRefresherCreateSecretSuccess(t *testing.T) {
	credProvider := &credentials.FakeDockerCredProvider{
		FnGetCreds: func() (credentials.DockerConfig, error) {
			return credentials.DockerConfig{}, nil
		},
		FnGetRefTime: func() time.Duration {
			return time.Minute * 1
		},
	}

	secretsClient := &k8s.FakeSecrets{
		FnCreate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, nil
		},
	}
	kubeClient := &k8s.FakeSecretsNamespacer{
		Fn: func(string) kcl.SecretsInterface {
			return secretsClient
		},
	}
	appAddedCh := make(chan string)
	tokenRefErrCh := make(chan error)
	doneCh := make(chan struct{})
	go tokenRefresher(kubeClient, credProvider, appAddedCh, tokenRefErrCh, doneCh)
	appAddedCh <- "testapp"
	timeoutCh := time.After(time.Second * 5)
	select {
	case err := <-tokenRefErrCh:
		assert.NoErr(t, err)
	case <-timeoutCh:
		close(doneCh)
	}
	assert.Equal(t, len(secretsClient.CreateCalls), 1, "number of create secret calls")
}

func TestTokenRefresherUpdateSecretSuccess(t *testing.T) {
	appList = []string{}
	credProvider := &credentials.FakeDockerCredProvider{
		FnGetCreds: func() (credentials.DockerConfig, error) {
			return credentials.DockerConfig{}, nil
		},
		FnGetRefTime: func() time.Duration {
			return time.Second * 2
		},
	}

	secretsClient := &k8s.FakeSecrets{
		FnCreate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, nil
		},
		FnUpdate: func(*api.Secret) (*api.Secret, error) {
			return &api.Secret{}, nil
		},
	}
	kubeClient := &k8s.FakeSecretsNamespacer{
		Fn: func(string) kcl.SecretsInterface {
			return secretsClient
		},
	}
	appAddedCh := make(chan string)
	tokenRefErrCh := make(chan error)
	doneCh := make(chan struct{})
	go tokenRefresher(kubeClient, credProvider, appAddedCh, tokenRefErrCh, doneCh)
	appList = append(appList, "testapp")
	appAddedCh <- "testapp"
	timeoutCh := time.After(time.Second * 5)
	select {
	case err := <-tokenRefErrCh:
		assert.NoErr(t, err)
	case <-timeoutCh:
		close(doneCh)
	}
	assert.Equal(t, len(secretsClient.CreateCalls), 1, "number of create secret calls")
	assert.Equal(t, len(secretsClient.UpdateCalls), 2, "number of update secret calls")
}
