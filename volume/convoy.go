package volume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
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

func (client *ConvoyClient) DeleteVolume(uuid string) error {
	reqBody, err := json.Marshal(api.VolumeDeleteRequest{
		VolumeUUID: uuid,
	})
	if err != nil {
		return err
	}

	err = client.doRequest("DELETE", "/v1/volumes/", reqBody, nil)
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

func (client *ConvoyClient) GetUUID(name string) (string, error) {
	v := url.Values{}
	v.Set(api.KEY_NAME, name)

	uuidResponse := &api.UUIDResponse{}
	err := client.doRequest("GET", "/uuid?"+v.Encode(), nil, uuidResponse)
	if err != nil {
		return "", err
	}

	return uuidResponse.UUID, nil
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

		return fmt.Errorf("Bad response: Code: %v. Response body: %s", resp.StatusCode, respBody)
	}

	if respTarget != nil {
		if err := json.NewDecoder(resp.Body).Decode(respTarget); err != nil {
			return err
		}
	}

	return nil
}
