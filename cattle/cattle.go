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
	rancherClient   *client.RancherClient
	driver          string
	storagepoolName string
}

func NewCattleClient(cattleUrl, cattleAccessKey, cattleSecretKey, spDriver, spName string) (*CattleClient, error) {
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
		rancherClient:   apiClient,
		driver:          spDriver,
		storagepoolName: spName,
	}, nil
}

func (c *CattleClient) CreateVolume(storagepoolUuid string, vol api.VolumeResponse) error {
	log.Debugf("create event %s", vol.UUID)
	eveResource := c.processVolume("volume.create", storagepoolUuid, vol)
	_, err := c.rancherClient.ExternalVolumeEvent.Create(eveResource)
	return err
}

func (c *CattleClient) processVolume(event, storagepoolUuid string, vol api.VolumeResponse) *client.ExternalVolumeEvent {
	opts := map[string]interface{}{}
	volume := client.Volume{
		Name:       vol.Name,
		Driver:     c.driver,
		DriverOpts: opts,
		ExternalId: vol.Name,
	}
	return &client.ExternalVolumeEvent{
		EventType:             event,
		ExternalId:            vol.Name,
		StoragePoolExternalId: storagepoolUuid,
		Volume:                volume,
	}
}

func (c *CattleClient) DeleteVolume(storagepoolUuid string, vol api.VolumeResponse) error {
	log.Debugf("delete event %s", vol.UUID)
	eveResource := c.processVolume("volume.delete", storagepoolUuid, vol)
	err := c.rancherClient.ExternalVolumeEvent.Delete(eveResource)
	return err
}

func (c *CattleClient) SyncStoragePool(storagepoolUuid string, hostUuids []string) error {
	log.Debugf("storagepool event %v", hostUuids)
	sp := client.StoragePool{
		Name:       c.storagepoolName,
		ExternalId: storagepoolUuid,
	}
	espe := &client.ExternalStoragePoolEvent{
		EventType:   "storagepool.create",
		HostUuids:   hostUuids,
		ExternalId:  storagepoolUuid,
		StoragePool: sp,
	}
	_, err := c.rancherClient.ExternalStoragePoolEvent.Create(espe)
	return err
}
