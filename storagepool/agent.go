package storagepool

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/convoy-agent/cattle"
)

var rootUuidFileName = "UUID"

type StoragepoolAgent struct {
	healthCheckInterval int
	storagepoolRootDir  string
	driver              string
	cattleClient        cattle.CattleInterface
}

func NewStoragepoolAgent(healthCheckInterval int, storagepoolRootDir, driver string, cattleClient cattle.CattleInterface) *StoragepoolAgent {
	return &StoragepoolAgent{
		healthCheckInterval: healthCheckInterval,
		storagepoolRootDir:  storagepoolRootDir,
		driver:              driver,
		cattleClient:        cattleClient,
	}
}

func (s *StoragepoolAgent) Run(metadataUrl string) error {
	prevSent := map[string]bool{}

	hc, err := newHealthChecker(metadataUrl)
	if err != nil {
		log.Errorf("Error initializing health checker, err = [%v]", err)
		return err
	}

	for {
		time.Sleep(time.Duration(s.healthCheckInterval) * time.Millisecond)

		currHosts, err := hc.populateHostMap()
		if err != nil {
			log.Errorf("Error while reading host info [%v]", err)
			continue
		}

		toSend := map[string]bool{}
		for uuid := range currHosts {
			toSend[uuid] = true
		}

		shouldSend := false
		for key := range toSend {
			if _, ok := prevSent[key]; !ok {
				shouldSend = true
			}
		}

		for key := range prevSent {
			if _, ok := toSend[key]; !ok {
				shouldSend = true
			}
		}

		if shouldSend {
			toSendList := []string{}
			for k := range toSend {
				toSendList = append(toSendList, k)
			}
			err := s.cattleClient.SyncStoragePool(s.driver, toSendList)
			if err != nil {
				log.Errorf("Error syncing storage pool events [%v]", err)
				continue
			}
			prevSent = toSend
		}
	}
	return nil
}
