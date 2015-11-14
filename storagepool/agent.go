package storagepool

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/convoy-agent/cattle"
)

var rootUuidFileName = "UUID"

type StoragepoolAgent struct {
	healthCheckInterval int
	storagepoolRootDir  string
	driver              string
	healthCheckBaseDir  string
	healthCheckType     string
	cattleClient        cattle.CattleInterface
}

func NewStoragepoolAgent(healthCheckInterval int, storagepoolRootDir, driver, healthCheckBaseDir, healthCheckType string, cattleClient cattle.CattleInterface) *StoragepoolAgent {
	return &StoragepoolAgent{
		healthCheckInterval: healthCheckInterval,
		storagepoolRootDir:  storagepoolRootDir,
		driver:              driver,
		healthCheckBaseDir:  healthCheckBaseDir,
		healthCheckType:     healthCheckType,
		cattleClient:        cattleClient,
	}
}

func (s *StoragepoolAgent) Run(metadataUrl string) error {
	// TODO Can we delete his non metadata healthcheck code?
	if _, err := os.Stat(filepath.Join(s.storagepoolRootDir, rootUuidFileName)); os.IsNotExist(err) && s.healthCheckType != "metadata" {
		err := ioutil.WriteFile(filepath.Join(s.storagepoolRootDir, rootUuidFileName), []byte(s.driver), 0644)
		if err != nil {
			return err
		}
	}

	prevHosts := map[string]string{}
	staleHosts := map[string]int{}
	prevSent := map[string]bool{}

	hc, err := newHealthChecker(s.healthCheckBaseDir, s.healthCheckType, metadataUrl)
	if err != nil {
		log.Errorf("Error initializing health checker, err = [%v]", err)
		return err
	}

	for {
		toDelete := []string{}
		time.Sleep(time.Duration(s.healthCheckInterval) * time.Second)

		currHosts, err := hc.populateHostMap()
		if err != nil {
			log.Errorf("Error while reading host info [%v]", err)
			continue
		}

		toSend := map[string]bool{}
		for uuid, stamp := range currHosts {
			prevStamp, ok := prevHosts[uuid]
			if !ok {
				toSend[uuid] = true
				continue
			}
			if s.healthCheckType != "metadata" && prevStamp == stamp {
				//stalehost
				staleHosts[uuid] = staleHosts[uuid] + 1
				if staleHosts[uuid] >= 3 {
					toDelete = append(toDelete, uuid)
					continue
				}
			}
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
			for _, uuid := range toDelete {
				if err := hc.deleteHost(uuid); err != nil {
					log.Error("error while deleting file [%v]", err)
				}
				delete(currHosts, uuid)
				delete(staleHosts, uuid)
			}
			prevSent = toSend
		}
		prevHosts = currHosts
	}
	return nil
}
