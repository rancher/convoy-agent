package storagepool

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type fileBasedHealthCheck struct {
	healthCheckBaseDir string
}

func (fc *fileBasedHealthCheck) populateHostMap() (map[string]string, error) {
	activeHosts := map[string]string{}
	fi, err := ioutil.ReadDir(fc.healthCheckBaseDir)
	if err != nil {
		return nil, err
	}
	for _, i := range fi {
		stamp, err := ioutil.ReadFile(filepath.Join(fc.healthCheckBaseDir, i.Name()))
		if err != nil {
			log.Errorf("Error reading file [%s] err [%v]", i.Name(), err)
			continue
		}
		activeHosts[i.Name()] = string(stamp)
	}
	return activeHosts, nil
}

func (fc *fileBasedHealthCheck) deleteHost(uuid string) error {
	return os.Remove(filepath.Join(fc.healthCheckBaseDir, uuid))
}
