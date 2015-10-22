package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
)

var (
	version    = 0
	selfStack  = metadata.Stack{}
	services   = []metadata.Service{}
	containers = []metadata.Container{}

	metadataUrl = ":12345"
)

func startMetadataServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/07-25-2015/version", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%d", version)
	})
	mux.HandleFunc("/07-25-2015/self/stack", func(w http.ResponseWriter, req *http.Request) {
		stackString, err := json.Marshal(selfStack)
		if err != nil {
			http.Error(w, "Could not marshall self/stack", 500)
			return
		}
		log.Error(string(stackString))
		fmt.Fprintf(w, string(stackString))
	})
	mux.HandleFunc("/07-25-2015/services", func(w http.ResponseWriter, req *http.Request) {
		servicesString, err := json.Marshal(services)
		if err != nil {
			http.Error(w, "Could not marshall services", 500)
			return
		}
		fmt.Fprintf(w, string(servicesString))
	})
	mux.HandleFunc("/07-25-2015/containers", func(w http.ResponseWriter, req *http.Request) {
		containerString, err := json.Marshal(containers)
		if err != nil {
			http.Error(w, "Could not marshall containers", 500)
			return
		}
		fmt.Fprintf(w, string(containerString))
	})
	if err := http.ListenAndServe(metadataUrl, mux); err != nil {
		log.Fatalf("error starting server err = [%v]", err)
	}
}

func setSelfStack(stack metadata.Stack) {
	updateVersion()
	selfStack = stack
}

func setContainers(conts []metadata.Container) {
	updateVersion()
	containers = conts
}

func setServices(servs []metadata.Service) {
	updateVersion()
	services = servs
}

func generateNewVersion() string {
	updateVersion()
	return strconv.Itoa(version)
}

func updateVersion() {
	version += 1
}
