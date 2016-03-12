package volume

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/convoy-agent/cattle"
)

type VolumeAgent struct {
	socketFile          string
	healthCheckInterval int
	volumeQueryInterval int
	cattleClient        cattle.CattleInterface
	driver              string
}

func NewVolumeAgent(socketFile string, volumeQueryInterval int, cattleClient cattle.CattleInterface, driver string) *VolumeAgent {
	return &VolumeAgent{
		socketFile:          socketFile,
		volumeQueryInterval: volumeQueryInterval,
		cattleClient:        cattleClient,
		driver:              driver,
	}
}

func (v *VolumeAgent) Run(controlChan chan bool) error {
	convoy, err := NewConvoyClient(v.socketFile)
	if err != nil {
		return err
	}

	vols := Volume{}

	for {
		select {
		case <-controlChan:
			controlChan <- true
			return nil
		case <-time.After(time.Duration(v.volumeQueryInterval) * time.Millisecond):
		}

		currVols, err := convoy.GetCurrVolumes()
		if err != nil {
			log.Error(err)
			continue
		}
		deletedVols := findDeletedVolumes(currVols, vols)
		createdVols := findCreatedVolumes(currVols, vols)

		for _, vol := range deletedVols {
			err := v.cattleClient.DeleteVolume(v.driver, vol)
			if err != nil {
				log.Errorf("Error sending delete event for volume name=[%s] err=[%v]", vol.Name, err)
				currVols[vol.Name] = vols[vol.Name]
			}
		}

		for _, vol := range createdVols {
			err := v.cattleClient.CreateVolume(v.driver, vol)
			if err != nil {
				log.Errorf("Error sending create event for volume name=[%s] err=[%v]", vol.Name, err)
				delete(currVols, vol.Name)
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
