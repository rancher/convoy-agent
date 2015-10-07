package volume

import (
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
