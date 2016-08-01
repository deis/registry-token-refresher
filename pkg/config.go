package pkg

import (
	"github.com/deis/registry-token-refresher/pkg/credentials"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	gcsKey          = "key.json"
	tokenRefreshKey = "DEIS_TOKEN_REFRESH_TIME"
)

func GetRegistryParams(registryCredLocation string) (credentials.Parameters, error) {
	params := make(map[string]interface{})
	files, err := ioutil.ReadDir(registryCredLocation)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || file.Name() == "..data" {
			continue
		}
		data, err := ioutil.ReadFile(registryCredLocation + file.Name())
		if err != nil {
			return nil, err
		}
		//GCS expects them to have the location of the service account credential json file
		if file.Name() == gcsKey {
			params["keyfile"] = registryCredLocation + file.Name()
		} else {
			params[file.Name()] = string(data)
		}
	}

	defaultRefreshTime := os.Getenv("tokenRefreshKey")
	if defaultRefreshTime != "" {
		refreshTime, err := strconv.ParseInt(defaultRefreshTime, 10, 32)
		if err != nil {
			return nil, err
		}
		params["defaultRefreshTime"] = refreshTime
	}

	return params, nil
}
