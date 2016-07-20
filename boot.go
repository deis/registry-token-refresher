package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/deis/registry-token-refresher/pkg"
	"github.com/deis/registry-token-refresher/pkg/credentials"
	"k8s.io/kubernetes/pkg/api"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
)

var appList []string
var registryLocation = os.Getenv("DEIS_REGISTRY_LOCATION")
var namespaceRefreshTime = os.Getenv("DEIS_NAMESPACE_REFRESH_TIME")

func getDiff(namespaceList []api.Namespace) []string {
	var added []string
	apps := make([]string, len(namespaceList))
	// create a set of app names
	appsSet := make(map[string]struct{}, len(appList))
	for _, app := range appList {
		appsSet[app] = struct{}{}
	}

	for _, ns := range namespaceList {
		if _, ok := appsSet[ns.Name]; !ok {
			added = append(added, ns.Name)
		}
		apps = append(apps, ns.Name)
	}

	appList = apps
	return added
}

func tokenRefresher(client kcl.SecretsNamespacer, credsProvider credentials.DockerCredProvider, appCh <-chan string, errCh chan<- error, doneCh <-chan struct{}) {
	creds, err := credsProvider.GetDockerCredentials()
	if err != nil {
		errCh <- err
		return
	}
	ticker := time.NewTicker(credsProvider.GetRefreshTime())
	defer ticker.Stop()
	for {
		select {
		case app := <-appCh:
			log.Printf("creating secret for app %s", app)
			if err = pkg.CreateSecret(client.Secrets(app), creds); err != nil {
				errCh <- err
				return
			}
		case <-ticker.C:
			creds, err = credsProvider.GetDockerCredentials()
			if err != nil {
				errCh <- err
				return
			}
			for _, app := range appList {
				log.Printf("updating secret for app %s", app)
				if err = pkg.UpdateSecret(client.Secrets(app), creds); err != nil {
					errCh <- err
					return
				}
			}
		case <-doneCh:
			return
		}
	}
}

func main() {
	kubeClient, err := kcl.NewInCluster()
	if err != nil {
		log.Fatal("Error getting kubernetes client ", err)
	}
	refreshTime, err := strconv.ParseInt(namespaceRefreshTime, 10, 32)
	if err != nil {
		log.Fatal("Error getting the namespace refresh time", err)
	}
	params, err := pkg.GetRegistryParams()
	if err != nil {
		log.Fatal("Error getting registry location credentials details", err)
	}
	credProvider, err := credentials.GetDockerCredentialsProvider(registryLocation, params)
	if err != nil {
		log.Fatal("Error getting credentials provider", err)
	}

	appAddedCh := make(chan string)
	tokenRefErrCh := make(chan error)
	doneCh := make(chan struct{})
	defer close(doneCh)

	go tokenRefresher(kubeClient, credProvider, appAddedCh, tokenRefErrCh, doneCh)

	for {
		select {
		case err = <-tokenRefErrCh:
			log.Fatalf("error during token refresh %s", err)
		default:
			labelMap := labels.Set{"heritage": "deis"}
			nsList, err := kubeClient.Namespaces().List(api.ListOptions{LabelSelector: labelMap.AsSelector(), FieldSelector: fields.Everything()})
			if err != nil {
				log.Fatal("Error getting kubernetes namespaces ", err)
			}
			added := getDiff(nsList.Items)
			for _, app := range added {
				appAddedCh <- app
			}
		}
		time.Sleep(time.Second * time.Duration(refreshTime))
	}
}
