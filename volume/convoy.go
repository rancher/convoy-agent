package volume

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/rancher/convoy/api"
)

type ConvoyClient struct {
	Client *http.Client
	Addr   string
}

type Volume map[string]api.VolumeResponse

func NewConvoyClient(socketFile string) (*ConvoyClient, error) {
	tr := &http.Transport{
		DisableCompression: true,
		Dial: func(_, _ string) (net.Conn, error) {
			return net.DialTimeout("unix", socketFile, 10*time.Second)
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	return &ConvoyClient{
		Client: client,
		Addr:   socketFile,
	}, nil
}

func (client *ConvoyClient) GetCurrVolumes() (Volume, error) {
	req, err := http.NewRequest("GET", "/v1/volumes/list", nil)
	if err != nil {
		return nil, err
	}

	req.URL.Host = client.Addr
	req.URL.Scheme = "http"

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	volumes := Volume{}
	err = json.Unmarshal(body, &volumes)
	return volumes, err
}

func (client *ConvoyClient) DeleteVolume(name string) error {
	reqBody, err := json.Marshal(api.VolumeDeleteRequest{
		VolumeName: name,
	})
	if err != nil {
		return err
	}

	err = client.doRequest("DELETE", "/v1/volumes/", reqBody, nil)
	if isNotFoundError(err) {
		return nil
	}
	return err
}

// Not really production worthy as it does not support driver options
func (client *ConvoyClient) CreateVolume(name string) error {
	reqBody, err := json.Marshal(api.VolumeCreateRequest{
		Name: name,
	})
	if err != nil {
		return err
	}

	err = client.doRequest("POST", "/v1/volumes/create", reqBody, nil)
	return err
}

func (client *ConvoyClient) GetVolume(name string) (*api.VolumeResponse, error) {
	reqBody, err := json.Marshal(api.VolumeInspectRequest{
		VolumeName: name,
	})
	if err != nil {
		return nil, err
	}

	vol := &api.VolumeResponse{}
	err = client.doRequest("GET", "/v1/volumes/", reqBody, vol)
	if isNotFoundError(err) {
		return nil, nil
	}
	return vol, err
}

func (client *ConvoyClient) doRequest(method string, path string, body []byte, respTarget interface{}) error {
	bodyBuf := bytes.NewBuffer(nil)
	if _, err := bodyBuf.Write(body); err != nil {
		return err
	}

	req, err := http.NewRequest(method, path, bodyBuf)
	if err != nil {
		return err
	}
	req.URL.Host = client.Addr
	req.URL.Scheme = "http"
	req.Header.Add("Context-Type", "application/json")

	resp, err := client.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return APIResponseError{
			ErrorMessage: string(respBody),
			StatusCode:   resp.StatusCode,
		}
	}

	if respTarget != nil {
		if err := json.NewDecoder(resp.Body).Decode(respTarget); err != nil {
			return err
		}
	}

	return nil
}

type APIResponseError struct {
	ErrorMessage string
	StatusCode   int
}

func (e APIResponseError) Error() string {
	return e.ErrorMessage
}

func isNotFoundError(err error) bool {
	if apiErr, ok := err.(APIResponseError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}
