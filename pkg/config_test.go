package pkg

import (
	"io/ioutil"
	"os"
	"os/user"
	"testing"
)

func TestGetregistryParams(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Logf("could not retrieve current user: %v", err)
		t.SkipNow()
	}
	if usr.Uid != "0" {
		t.Logf("current user does not have UID of zero (got %s) "+
			"so cannot create registry cred location, skipping", usr.Uid)
		t.SkipNow()
	}

	if err = os.MkdirAll(registryCredLocation, os.ModeDir); err != nil {
		t.Fatalf("could not create registry cred location: %v", err)
	}

	// start by writing out a file to registryCredLocation
	data := []byte("hello world\n")
	if err = ioutil.WriteFile(registryCredLocation+"foo", data, 0644); err != nil {
		t.Fatalf("could not write file to registry cred location: %v", err)
	}

	params, err := GetRegistryParams()
	if err != nil {
		t.Errorf("received error while retrieving registry params: %v", err)
	}

	val, ok := params["foo"]
	if !ok {
		t.Error("key foo does not exist in registry params")
	}
	if val != string(data) {
		t.Errorf("expected: %s got: %s", string(data), val)
	}

	// create a directory inside registry cred location, expecting it to pass
	if err = os.Mkdir(registryCredLocation+"bar", os.ModeDir); err != nil {
		t.Fatalf("could not create dir %s: %v", registryCredLocation+"bar", err)
	}

	_, err = GetRegistryParams()
	if err != nil {
		t.Errorf("received error while retrieving registry params: %v", err)
	}

	// create the special "..data" directory symlink, expecting it to pass
	if err = os.Symlink(registryCredLocation+"bar", registryCredLocation+"..data"); err != nil {
		t.Fatalf("could not create dir symlink ..data -> %s: %v", registryCredLocation+"bar", err)
	}

	_, err = GetRegistryParams()
	if err != nil {
		t.Errorf("received error while retrieving registry params: %v", err)
	}
}
