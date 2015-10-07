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
	storagepoolUuid     string
	healthCheckBaseDir  string
	cattleClient        cattle.CattleInterface
}

func NewStoragepoolAgent(healthCheckInterval int, storagepoolRootDir, storagepoolUuid, healthCheckBaseDir string, cattleClient cattle.CattleInterface) *StoragepoolAgent {
	return &StoragepoolAgent{
		healthCheckInterval: healthCheckInterval,
		storagepoolRootDir:  storagepoolRootDir,
		storagepoolUuid:     storagepoolUuid,
		healthCheckBaseDir:  healthCheckBaseDir,
		cattleClient:        cattleClient,
	}
}

func (s *StoragepoolAgent) Run() error {
	if _, err := os.Stat(filepath.Join(s.storagepoolRootDir, rootUuidFileName)); os.IsNotExist(err) {
		err := ioutil.WriteFile(filepath.Join(s.storagepoolRootDir, rootUuidFileName), []byte(s.storagepoolUuid), 0644)
		if err != nil {
			return err
		}
	}

	prevHosts := map[string]string{}
	staleHosts := map[string]int{}
	prevSent := map[string]bool{}

	for {
		toDelete := []string{}
		time.Sleep(time.Duration(s.healthCheckInterval) * time.Second)

		currHosts, err := populateHostMap(s.healthCheckBaseDir)
		if err != nil {
			log.Error("Error while reading host info [%v]", err)
			continue
		}

		toSend := map[string]bool{}
		for uuid, stamp := range currHosts {
			prevStamp, ok := prevHosts[uuid]
			if !ok {
				toSend[uuid] = true
				continue
			}
			if prevStamp == stamp {
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
			err := s.cattleClient.SyncStoragePool(s.storagepoolUuid, toSendList)
			if err != nil {
				log.Errorf("Error syncing storage pool events [%v]", err)
				continue
			}
			for _, uuid := range toDelete {
				if err := deleteHost(uuid, s.healthCheckBaseDir); err != nil {
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

func deleteHost(uuid, healthCheckBaseDir string) error {
	return os.Remove(filepath.Join(healthCheckBaseDir, uuid))
}

func populateHostMap(healthCheckBaseDir string) (map[string]string, error) {
	activeHosts := map[string]string{}
	f, err := os.Open(healthCheckBaseDir)
	if err != nil {
		return nil, err
	}
	fi, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, i := range fi {
		stamp, err := ioutil.ReadFile(filepath.Join(healthCheckBaseDir, i.Name()))
		if err != nil {
			log.Errorf("Error reading file [%s] err [%v]", i.Name(), err)
			continue
		}
		activeHosts[i.Name()] = string(stamp)
	}
	return activeHosts, nil
}
