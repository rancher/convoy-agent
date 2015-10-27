package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/convoy-agent/storagepool"
)

var (
	sphc       = ".healthcheck2"
	hcFileType = "file"
)

func writeTestHostHealthcheckFile(filename string, controlChan chan bool) {
	for {
		err := ioutil.WriteFile(filepath.Join(sphc, filename), []byte(time.Now().Format(time.RFC1123Z)), 0644)
		if err != nil {
			log.Errorf("Error writing to file [%s] err [%v]", filename, err)
			continue
		}
		select {
		case v := <-controlChan:
			if v {
				//if true, we pause until a resume event is sent
				<-controlChan
			} else {
				os.Remove(filepath.Join(hc, filename))
				// else, it is an exit event and we exit
				return
			}
		case <-time.After(5 * time.Second):
		}
	}
}

func TestSyncInitialHostsInStoragePool(t *testing.T) {
	controlChan1 := make(chan bool, 1)
	controlChan2 := make(chan bool, 1)
	go writeTestHostHealthcheckFile("hostUuid1", controlChan1)
	go writeTestHostHealthcheckFile("hostUuid2", controlChan2)

	tc := &testCattleClient{}

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcFileType, tc)
		err := spAgent.Run("")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()

	<-time.After(10 * time.Second)

	uuids := tc.getLastSync()
	if len(uuids) != 2 {
		t.Fatalf("expected 2 storagepool sync events, but received [%d] %v", len(uuids), uuids)
	}

	uuid1found := false
	uuid2found := false

	for _, uuid := range uuids {
		if uuid == "hostUuid1" {
			uuid1found = true
		}
		if uuid == "hostUuid2" {
			uuid2found = true
		}
	}

	if !uuid1found || !uuid2found {
		t.Fatalf("sync event not as expected, received %v", uuids)
	}

	controlChan1 <- false
	controlChan2 <- false
}

func TestSyncLostHostDetectedInStoragePool(t *testing.T) {
	controlChan1 := make(chan bool, 1)
	controlChan2 := make(chan bool, 1)
	go writeTestHostHealthcheckFile("hostUuid1", controlChan1)
	go writeTestHostHealthcheckFile("hostUuid2", controlChan2)

	tc := &testCattleClient{}

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcFileType, tc)
		err := spAgent.Run("")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()

	<-time.After(10 * time.Second)

	uuids := tc.getLastSync()
	if len(uuids) != 2 {
		t.Fatalf("expected 2 hosts in storagepool sync event, but received %d", len(uuids))
	}

	uuid1found := false
	uuid2found := false

	for _, uuid := range uuids {
		if uuid == "hostUuid1" {
			uuid1found = true
		}
		if uuid == "hostUuid2" {
			uuid2found = true
		}
	}

	if !uuid1found || !uuid2found {
		t.Fatalf("sync event not as expected, received %v", uuids)
	}

	// pause first host
	controlChan1 <- true

	<-time.After(20 * time.Second)

	uuids = tc.getLastSync()
	if len(uuids) != 1 {
		t.Fatalf("expected 1 host in storagepool sync event, but received %s", len(uuids))
	}

	if uuids[0] != "hostUuid2" {
		t.Fatalf("Expected hostUuid2 to be present,but found %s", uuids[0])
	}

	// resume
	controlChan1 <- true

	controlChan1 <- false
	controlChan2 <- false
}

func TestSyncNewHostDetectedInStoragePool(t *testing.T) {
	controlChan1 := make(chan bool, 1)
	controlChan2 := make(chan bool, 1)
	go writeTestHostHealthcheckFile("hostUuid1", controlChan1)
	go writeTestHostHealthcheckFile("hostUuid2", controlChan2)

	tc := &testCattleClient{}

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcFileType, tc)
		err := spAgent.Run("")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()

	<-time.After(10 * time.Second)

	uuids := tc.getLastSync()
	if len(uuids) != 2 {
		t.Fatalf("expected 2 storagepool sync events, but received %d", len(uuids))
	}

	uuid1found := false
	uuid2found := false
	uuid3found := false

	for _, uuid := range uuids {
		if uuid == "hostUuid1" {
			uuid1found = true
		}
		if uuid == "hostUuid2" {
			uuid2found = true
		}
	}

	if !uuid1found || !uuid2found {
		t.Fatalf("sync event not as expected, received %v", uuids)
	}
	controlChan3 := make(chan bool, 1)
	go writeTestHostHealthcheckFile("hostUuid3", controlChan3)

	<-time.After(10 * time.Second)

	uuids = tc.getLastSync()
	if len(uuids) != 3 {
		t.Fatalf("expected 3 host in storagepool sync event, but received %s", len(uuids))
	}

	for _, uuid := range uuids {
		if uuid == "hostUuid1" {
			uuid1found = true
		}
		if uuid == "hostUuid2" {
			uuid2found = true
		}
		if uuid == "hostUuid3" {
			uuid3found = true
		}
	}

	if !uuid1found || !uuid2found || !uuid3found {
		t.Fatalf("sync event not as expected, received %v", uuids)
	}

	controlChan1 <- false
	controlChan2 <- false
	controlChan3 <- false
}
