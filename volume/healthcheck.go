package volume

import (
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
		f, err := os.OpenFile(filepath.Join(baseDir, hostUuid), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Errorf("error opening file to write err=[%v]", err)
			continue
		}
		_, err = f.Write([]byte(time.Now().Format(time.RFC1123Z)))
		if err != nil {
			log.Errorf("error writing healthcheck file err=[%v]", err)
			f.Close()
			continue
		}
		err = f.Sync()
		if err != nil {
			log.Errorf("error syncing file to filesystem err=[%v]", err)
			f.Close()
			continue
		}
		f.Close()
	}
}
