package volume

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

func writeHealthCheckFile(hostUuid, baseDir string, healthCheckInterval int, controlChan chan bool) {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.Mkdir(baseDir, 0644); err != nil {
			log.Error("Error creating baseDir for healthcheck.")
			controlChan <- true
		}
	}

	for {
		select {
		case <-controlChan:
			controlChan <- true
			return
		case <-time.After(5 * time.Second):
		}
		err := ioutil.WriteFile(filepath.Join(baseDir, hostUuid), []byte(time.Now().Format(time.RFC1123Z)), 0644)
		if err != nil {
			log.Error("error writing healthcheck file err=[%v]", err)
		}
	}
}
