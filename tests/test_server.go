package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rancher/convoy-agent/volume"
	"github.com/rancher/convoy/api"
)

var (
	socketFile = "/var/run/convoy/convoy.sock"
	convoy     *volume.ConvoyClient
)

func init() {
	var err error
	convoy, err = volume.NewConvoyClient(socketFile)
	if err != nil {
		panic(fmt.Sprintf("ERROR [%v]. Could not connect to convoy", err))
	}
}

func createVolume(name string) (string, error) {
	inp := &api.VolumeCreateRequest{
		Name: name,
	}
	reqString, err := json.Marshal(inp)
	if err != nil {
		return "", err
	}
	reqBody := bytes.NewBuffer(nil)
	if _, err := reqBody.Write(reqString); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "/v1/volumes/create", reqBody)
	if err != nil {
		return "", err
	}
	req.URL.Host = socketFile
	req.URL.Scheme = "http"

	resp, err := convoy.Client.Do(req)
	if err != nil {
		return "", err
	}
	respString, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", respString), nil
}

func deleteVolume(name string) (string, error) {
	inp := &api.VolumeDeleteRequest{
		VolumeName: name,
	}
	reqString, err := json.Marshal(inp)
	if err != nil {
		return "", err
	}
	reqBody := bytes.NewBuffer(nil)
	if _, err := reqBody.Write(reqString); err != nil {
		return "", err
	}

	req, err := http.NewRequest("DELETE", "/v1/volumes/", reqBody)
	if err != nil {
		return "", err
	}
	req.URL.Host = socketFile
	req.URL.Scheme = "http"

	resp, err := convoy.Client.Do(req)
	if err != nil {
		return "", err
	}
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", respString), nil
}
