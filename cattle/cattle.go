package cattle

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/convoy/api"
	"github.com/rancher/go-rancher/client"
)

type CattleInterface interface {
	CreateVolume(string, api.VolumeResponse) error
	DeleteVolume(string, api.VolumeResponse) error
	SyncStoragePool(string, []string) error
}

type CattleClient struct {
	rancherClient *client.RancherClient
}

func NewCattleClient(cattleUrl, cattleAccessKey, cattleSecretKey string) (*CattleClient, error) {
	if cattleUrl == "" {
		return nil, errors.New("cattle url is empty")
	}

	apiClient, err := client.NewRancherClient(&client.ClientOpts{
		Url:       cattleUrl,
		AccessKey: cattleAccessKey,
		SecretKey: cattleSecretKey,
	})

	if err != nil {
		return nil, err
	}

	return &CattleClient{
		rancherClient: apiClient,
	}, nil
}

func (c *CattleClient) CreateVolume(driver string, vol api.VolumeResponse) error {
	log.Debugf("create event %s", vol.Name)
	eveResource := c.processVolume("volume.create", driver, vol)
	_, err := c.rancherClient.ExternalVolumeEvent.Create(eveResource)
	return err
}

func (c *CattleClient) processVolume(event, driver string, vol api.VolumeResponse) *client.ExternalVolumeEvent {
	opts := map[string]interface{}{}
	volume := client.Volume{
		Name:       vol.Name,
		Driver:     driver,
		DriverOpts: opts,
		ExternalId: vol.Name,
	}
	return &client.ExternalVolumeEvent{
		EventType:  event,
		ExternalId: vol.Name,
		Volume:     volume,
	}
}

func (c *CattleClient) DeleteVolume(driver string, vol api.VolumeResponse) error {
	log.Debugf("delete event %s", vol.Name)
	eveResource := c.processVolume("volume.delete", driver, vol)
	_, err := c.rancherClient.ExternalVolumeEvent.Create(eveResource)
	return err
}

func (c *CattleClient) SyncStoragePool(driver string, hostUuids []string) error {
	log.Debugf("storagepool event %v", hostUuids)
	sp := client.StoragePool{
		Name:       driver,
		ExternalId: driver,
		DriverName: driver,
	}
	espe := &client.ExternalStoragePoolEvent{
		EventType:   "storagepool.create",
		HostUuids:   hostUuids,
		ExternalId:  driver,
		StoragePool: sp,
	}
	_, err := c.rancherClient.ExternalStoragePoolEvent.Create(espe)
	return err
}
