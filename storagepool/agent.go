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

	prevHosts, err := getActiveHosts(s.healthCheckBaseDir, s.healthCheckInterval)
	if err != nil {
		return err
	}
	prevSent := map[string]bool{}
	toSendList := []string{}
	for hostUuid := range prevHosts {
		prevSent[hostUuid] = true
		toSendList = append(toSendList, hostUuid)
	}
	s.cattleClient.SyncStoragePool(s.storagepoolUuid, toSendList)

	for {
		<-time.After(time.Duration(s.healthCheckInterval) * time.Second)

		currHosts, err := populateHostMap(s.healthCheckBaseDir)
		if err != nil {
			log.Error("Error while reading host info [%v]", err)
			continue
		}

		toSend := map[string]bool{}
		for uuid, timeStamp := range currHosts {
			if time.Now().Sub(timeStamp) > time.Duration(3*s.healthCheckInterval)*time.Second {
				if err := deleteHost(uuid, s.healthCheckBaseDir); err != nil {
					log.Error("error while deleting file [%v]", err)
				}
				continue
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
			s.cattleClient.SyncStoragePool(s.storagepoolUuid, toSendList)
			prevSent = toSend
		}
		prevHosts = currHosts
	}
	return nil
}

func deleteHost(uuid, healthCheckBaseDir string) error {
	return os.Remove(filepath.Join(healthCheckBaseDir, uuid))
}

func getActiveHosts(healthCheckBaseDir string, healthCheckInterval int) (map[string]time.Time, error) {
	refHostMap, err := populateHostMap(healthCheckBaseDir)
	if err != nil {
		return nil, err
	}

	<-time.After(time.Duration(healthCheckInterval) * time.Second)
	activeHosts, err := populateHostMap(healthCheckBaseDir)
	if err != nil {
		return nil, err
	}
	prevHosts := map[string]time.Time{}
	for uuid, newTime := range activeHosts {
		oldTime, ok := refHostMap[uuid]
		if !ok || newTime.Unix() > oldTime.Unix() {
			prevHosts[uuid] = newTime
		} else {
			//it is stale
			deleteHost(uuid, healthCheckBaseDir)
		}
	}
	return prevHosts, nil
}

func populateHostMap(healthCheckBaseDir string) (map[string]time.Time, error) {
	activeHosts := map[string]time.Time{}
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
		t, err := time.Parse(time.RFC1123Z, string(stamp))
		if err != nil {
			log.Errorf("Error reading time from healthcheck file [%s] err [%v]", i.Name(), err)
		}
		activeHosts[i.Name()] = t
	}
	return activeHosts, nil
}
