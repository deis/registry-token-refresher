package k8s

import (
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/watch"
)

// FakeSecrets is a mock function that can be swapped in for an SecretGetter or
// (k8s.io/kubernetes/pkg/client/unversioned).SecretsInterface, so you can
// unit test your code.
type FakeSecrets struct {
	FnCreate    func(*api.Secret) (*api.Secret, error)
	FnUpdate    func(*api.Secret) (*api.Secret, error)
	CreateCalls []*api.Secret
	UpdateCalls []*api.Secret
}

// Get is the interface definition.
func (f *FakeSecrets) Get(name string) (*api.Secret, error) {
	return &api.Secret{}, nil
}

// Delete is the interface definition.
func (f *FakeSecrets) Delete(name string) error {
	return nil
}

// Create is the interface definition.
func (f *FakeSecrets) Create(secret *api.Secret) (*api.Secret, error) {
	f.CreateCalls = append(f.CreateCalls, secret)
	return f.FnCreate(secret)
}

// Update is the interface definition.
func (f *FakeSecrets) Update(secret *api.Secret) (*api.Secret, error) {
	f.UpdateCalls = append(f.UpdateCalls, secret)
	return f.FnUpdate(secret)
}

// List is the interface definition.
func (f *FakeSecrets) List(opts api.ListOptions) (*api.SecretList, error) {
	return &api.SecretList{}, nil
}

// Watch is the interface definition.
func (f *FakeSecrets) Watch(opts api.ListOptions) (watch.Interface, error) {
	return nil, nil
}

// FakeSecretsNamespacer is a mock function that can be swapped in for an
// (k8s.io/kubernetes/pkg/client/unversioned).SecretsNamespacer, so you can unit test you code
type FakeSecretsNamespacer struct {
	Fn func(string) client.SecretsInterface
}

// Secrets is the interface definition.
func (f *FakeSecretsNamespacer) Secrets(namespace string) client.SecretsInterface {
	return f.Fn(namespace)
}
