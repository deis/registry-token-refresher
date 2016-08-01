package pkg

import (
	"encoding/json"
	"github.com/deis/registry-token-refresher/pkg/credentials"
	"k8s.io/kubernetes/pkg/api"
	apierrors "k8s.io/kubernetes/pkg/api/errors"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"os"
)

var secretName = os.Getenv("DEIS_REGISTRY_SECRET_PREFIX") + "-" + os.Getenv("DEIS_REGISTRY_LOCATION")

type AuthsStruct struct {
	Auths map[string]AuthInfo `json:"auths"`
}
type AuthInfo struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

// CreateSecret creates a imagepull secret if it doesn't exist or updates the existing one
func CreateSecret(secretsClient kcl.SecretsInterface, dockerconfig credentials.DockerConfig) error {
	authStruct := AuthsStruct{
		Auths: map[string]AuthInfo{
			dockerconfig.Hostname: AuthInfo{
				Auth:  dockerconfig.Token,
				Email: "none@none.com",
			},
		},
	}
	auth, err := json.Marshal(authStruct)
	if err != nil {
		return err
	}
	newSecret := new(api.Secret)
	newSecret.Name = secretName
	newSecret.Type = api.SecretTypeDockerConfigJson
	newSecret.Data = make(map[string][]byte)
	newSecret.Data[api.DockerConfigJsonKey] = auth
	if _, err := secretsClient.Create(newSecret); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return UpdateSecret(secretsClient, dockerconfig)
		}
		return err
	}
	return nil
}

// UpdateSecret update the imagepull secret if it exists
func UpdateSecret(secretsClient kcl.SecretsInterface, dockerconfig credentials.DockerConfig) error {
	authStruct := AuthsStruct{
		Auths: map[string]AuthInfo{
			dockerconfig.Hostname: AuthInfo{
				Auth:  dockerconfig.Token,
				Email: "none@none.com",
			},
		},
	}
	auth, err := json.Marshal(authStruct)
	if err != nil {
		return err
	}

	newSecret := new(api.Secret)
	newSecret.Name = secretName
	newSecret.Type = api.SecretTypeDockerConfigJson
	newSecret.Data = make(map[string][]byte)
	newSecret.Data[api.DockerConfigJsonKey] = auth
	if _, err := secretsClient.Update(newSecret); err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}
