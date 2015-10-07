package volume

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/convoy-agent/cattle"
)

type VolumeAgent struct {
	healthCheckBaseDir  string
	socketFile          string
	healthCheckInterval int
	cattleClient        cattle.CattleInterface
	storagepoolUuid     string
	hostUuid            string
}

func NewVolumeAgent(healthCheckBaseDir, socketFile, hostUuid string, healthCheckInterval int, cattleClient cattle.CattleInterface, storagepoolUuid string) *VolumeAgent {
	return &VolumeAgent{
		healthCheckBaseDir:  healthCheckBaseDir,
		socketFile:          socketFile,
		healthCheckInterval: healthCheckInterval,
		cattleClient:        cattleClient,
		storagepoolUuid:     storagepoolUuid,
		hostUuid:            hostUuid,
	}
}

func (v *VolumeAgent) Run(controlChan chan bool) error {
	convoy, err := NewConvoyClient(v.socketFile)
	if err != nil {
		return err
	}

	go writeHealthCheckFile(v.hostUuid, v.healthCheckBaseDir, v.healthCheckInterval, controlChan)

	vols := Volume{}

	for {
		select {
		case <-controlChan:
			controlChan <- true
			return nil
		case <-time.After(1 * time.Second):
		}

		currVols, err := convoy.GetCurrVolumes()
		if err != nil {
			log.Error(err)
			continue
		}
		deletedVols := findDeletedVolumes(currVols, vols)
		createdVols := findCreatedVolumes(currVols, vols)

		for _, vol := range deletedVols {
			err := v.cattleClient.DeleteVolume(v.storagepoolUuid, vol)
			if err != nil {
				log.Errorf("Error sending delete event for volume ID=[%s] err=[%v]", vol.UUID, err)
				currVols[vol.UUID] = vols[vol.UUID]
			}
		}

		for _, vol := range createdVols {
			err := v.cattleClient.CreateVolume(v.storagepoolUuid, vol)
			if err != nil {
				log.Errorf("Error sending create event for volume ID=[%s] err=[%v]", vol.UUID, err)
				delete(currVols, vol.UUID)
			}
		}
		vols = currVols
	}
	return nil
}

func findDeletedVolumes(curr, prev Volume) Volume {
	deleted := Volume{}
	for key, vol := range prev {
		if _, ok := curr[key]; !ok {
			deleted[key] = vol
		}
	}
	return deleted
}

func findCreatedVolumes(curr, prev Volume) Volume {
	created := Volume{}
	for key, vol := range curr {
		if _, ok := prev[key]; !ok {
			created[key] = vol
		}
	}
	return created
}
