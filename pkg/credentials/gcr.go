package credentials

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	gcrScope = "https://www.googleapis.com/auth/devstorage.read_write"
)

type gcrDockerCredProvider struct {
	params Parameters
}

func (d *gcrDockerCredProvider) GetDockerCredentials() (DockerConfig, error) {
	srvAccount, err := ioutil.ReadFile(fmt.Sprint(d.params["keyfile"]))
	if err != nil {
		return DockerConfig{}, err
	}
	config, err := google.JWTConfigFromJSON(srvAccount, gcrScope)
	if err != nil {
		return DockerConfig{}, err
	}
	ts := config.TokenSource(oauth2.NoContext)
	token, err := ts.Token()
	if err != nil {
		return DockerConfig{}, err
	}
	hostname := d.params["hostname"]
	if hostname == nil || fmt.Sprint(hostname) == "" {
		hostname = "https://gcr.io"
	}
	encToken := base64.StdEncoding.EncodeToString([]byte("oauth2accesstoken:" + token.AccessToken))
	dockerConfig := DockerConfig{Token: encToken,
		ExpiresAt: token.Expiry,
		Hostname:  fmt.Sprint(hostname)}
	return dockerConfig, nil
}

func (d *gcrDockerCredProvider) GetRefreshTime() time.Duration {
	defaultRefreshTime := d.params["defaultRefreshTime"]
	if defaultRefreshTime != nil {
		return time.Minute * time.Duration(defaultRefreshTime.(int))
	}
	return time.Minute * 25
}
