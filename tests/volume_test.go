package tests

import (
	"testing"
	"time"

	"github.com/rancher/convoy-agent/volume"
)

var (
	sp = "1234567890"
)

func waitForCreateDeleteEvents(t *testCattleClient, createChan chan<- string, deleteChan chan<- string, controlChan chan bool) {
	for {
		select {
		case <-controlChan:
			controlChan <- true
			return
		case <-time.After(100 * time.Millisecond):
		}
		ev := t.getLastEvent()
		if ev != "" {
			if len(ev) < len("CREATED_") {
				continue
			}
			if ev[:len("CREATED_")] == "CREATED_" {
				createChan <- ev[len("CREATED_"):]
			} else if ev[:len("DELETED_")] == "DELETED_" {
				deleteChan <- ev[len("DELETED_"):]
			}
		}
	}

}

func TestCreateEventsOnStartup(t *testing.T) {
	createChan := make(chan string, 10)
	deleteChan := make(chan string, 10)

	controlChan := make(chan bool, 1)

	tc := &testCattleClient{}

	uuid1, err := createVolume("testVol1")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err =[%v]", uuid1, err)
	}
	defer deleteVolume(uuid1)
	uuid2, err := createVolume("testVol2")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err =[%v]", uuid2, err)
	}
	defer deleteVolume(uuid2)
	uuid3, err := createVolume("testVol3")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err =[%v]", uuid2, err)
	}
	defer deleteVolume(uuid3)

	go func() {
		va := volume.NewVolumeAgent(socketFile, 100, tc, sp)
		err := va.Run(controlChan)
		if err != nil {
			t.Fatalf("Error starting convoy agent err=[%v]", err)
		}
	}()

	go waitForCreateDeleteEvents(tc, createChan, deleteChan, controlChan)

	obtainedUuids := []string{}

	for i := 0; i < 3; i++ {
		select {
		case vol := <-createChan:
			obtainedUuids = append(obtainedUuids, vol)
		case <-time.After(10 * time.Second):
			t.Fatal("All volume events were not received")
		}
	}

	found1 := true
	found2 := true
	found3 := true

	if len(obtainedUuids) != 3 {
		t.Errorf("created 3 vols but obtained events for %d vols", len(obtainedUuids))
		t.Fail()
	}

	for _, uuid := range obtainedUuids {
		if uuid == uuid1 {
			found1 = true
		}
		if uuid == uuid2 {
			found2 = true
		}
		if uuid == uuid3 {
			found3 = true
		}
	}

	if !found1 || !found2 || !found3 {
		t.Error("obtained Uuids do not match expected Uuids")
	}
	controlChan <- true
}

func TestNewCreateEventsAfterStartup(t *testing.T) {
	createChan := make(chan string, 10)
	deleteChan := make(chan string, 10)

	controlChan := make(chan bool, 1)

	tc := &testCattleClient{}

	go func() {
		va := volume.NewVolumeAgent(socketFile, 100, tc, sp)
		err := va.Run(controlChan)
		if err != nil {
			t.Fatalf("Error starting convoy agent err=[%v]", err)
		}
	}()
	uuid1, err := createVolume("testVol1")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err=[%v]", uuid1, err)
	}
	defer deleteVolume(uuid1)

	go waitForCreateDeleteEvents(tc, createChan, deleteChan, controlChan)

	select {
	case vol := <-createChan:
		if vol != uuid1 {
			t.Errorf("Excpected uuid = [%s] but obtained uuid = [%s]", uuid1, vol)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Volume create event not detected")
	}
	uuid2, err := createVolume("testVol2")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err=[%v]", uuid2, err)
	}
	defer deleteVolume(uuid2)

	select {
	case vol := <-createChan:
		if vol != uuid2 {
			t.Errorf("Excpected uuid = [%s] but obtained uuid = [%s]", uuid2, vol)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Volume create event not detected")
	}
	controlChan <- true
}

func TestNewDeleteEventsAfterStartup(t *testing.T) {
	createChan := make(chan string, 10)
	deleteChan := make(chan string, 10)

	controlChan := make(chan bool, 1)

	tc := &testCattleClient{}

	go func() {
		va := volume.NewVolumeAgent(socketFile, 100, tc, sp)
		err := va.Run(controlChan)
		if err != nil {
			t.Fatalf("Error starting convoy agent err=[%v]", err)
		}
	}()
	uuid1, err := createVolume("testVol1")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err=[%v]", uuid1, err)
	}

	go waitForCreateDeleteEvents(tc, createChan, deleteChan, controlChan)

	select {
	case vol := <-createChan:
		if vol != uuid1 {
			t.Errorf("Excpected uuid = [%s] but obtained uuid = [%s]", uuid1, vol)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Volume create event not detected")
	}
	uuid2, err := createVolume("testVol2")
	if err != nil {
		t.Fatalf("Error while creating volume UUID=[%s] err=[%v]", uuid2, err)
	}

	select {
	case vol := <-createChan:
		if vol != uuid2 {
			t.Errorf("Excpected uuid = [%s] but obtained uuid = [%s]", uuid2, vol)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Volume create event not detected")
	}

	deleteVolume(uuid2)

	select {
	case vol := <-deleteChan:
		if vol != uuid2 {
			t.Errorf("Excpected uuid = [%s] but obtained uuid = [%s]", uuid2, vol)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Volume delete event not detected")
	}

	deleteVolume(uuid1)

	select {
	case vol := <-deleteChan:
		if vol != uuid1 {
			t.Errorf("Excpected uuid = [%s] but obtained uuid = [%s]", uuid1, vol)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Volume delete event not detected")
	}
	controlChan <- true

}
