package tests

import (
	"os"
	"testing"
	"time"

	"github.com/rancher/convoy-agent/storagepool"
	"github.com/rancher/go-rancher-metadata/metadata"
)

var (
	hcMetadataType = "metadata"
)

func TestMain(m *testing.M) {
	go startMetadataServer()
	os.Exit(m.Run())
}

func TestReadsHostsCorrectly(t *testing.T) {
	tc := &testCattleClient{}

	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service2",
		},
	}
	setSelfStack(stack)

	services := []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	setServices(services)

	containers := []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container2",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid2",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
		{
			Name:        "container5",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid5",
		},
		{
			Name:        "container6",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid6",
		},
	}
	setContainers(containers)

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcMetadataType, tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/07-25-2015")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(10 * time.Second)

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

}

func TestDetectsVersionChange(t *testing.T) {
	tc := &testCattleClient{}
	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service2",
		},
	}
	setSelfStack(stack)

	services := []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	setServices(services)

	containers := []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container2",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid2",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
		{
			Name:        "container5",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid5",
		},
		{
			Name:        "container6",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid6",
		},
	}
	setContainers(containers)
	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcMetadataType, tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/07-25-2015")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(10 * time.Second)

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

	stack = metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service3",
		},
	}
	setSelfStack(stack)

	services = []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	setServices(services)

	containers = []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container2",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid2",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
		{
			Name:        "container5",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid5",
		},
		{
			Name:        "container6",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid6",
		},
	}
	setContainers(containers)
	time.Sleep(10 * time.Second)

	uuids = tc.getLastSync()
	if len(uuids) != 4 {
		t.Fatalf("expected 4 storagepool sync events, but received [%d] %v", len(uuids), uuids)
	}

	uuid1found = false
	uuid2found = false
	uuid3found := false
	uuid4found := false

	for _, uuid := range uuids {
		if uuid == "hostUuid1" {
			uuid1found = true
		}
		if uuid == "hostUuid2" {
			uuid2found = true
		}
		if uuid == "hostUuid5" {
			uuid3found = true
		}
		if uuid == "hostUuid6" {
			uuid4found = true
		}
	}

	if !uuid1found || !uuid2found || !uuid3found || !uuid4found {
		t.Fatalf("sync event not as expected, received %v", uuids)
	}
}

func TestNoVersionChange(t *testing.T) {
	tc := &testCattleClient{}
	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service2",
		},
	}
	setSelfStack(stack)

	services := []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	setServices(services)

	containers := []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container2",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid2",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
		{
			Name:        "container5",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid5",
		},
		{
			Name:        "container6",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid6",
		},
	}
	setContainers(containers)
	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcMetadataType, tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/07-25-2015")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(10 * time.Second)

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

	stack = metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service3",
		},
	}
	selfStack = stack

	newServices := []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	services = newServices

	newContainers := []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container2",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid2",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
		{
			Name:        "container5",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid5",
		},
		{
			Name:        "container6",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid6",
		},
	}
	containers = newContainers
	time.Sleep(10 * time.Second)

	uuids = tc.getLastSync()
	if len(uuids) != 0 {
		t.Fatalf("expected 0 storagepool sync events, but received [%d] %v", len(uuids), uuids)
	}

}

func TestVersionChangeAndDeletion(t *testing.T) {
	tc := &testCattleClient{}
	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service2",
		},
	}
	setSelfStack(stack)

	services := []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	setServices(services)

	containers := []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container2",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid2",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
		{
			Name:        "container5",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid5",
		},
		{
			Name:        "container6",
			ServiceName: "service3",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid6",
		},
	}
	setContainers(containers)
	go func() {
		spAgent := storagepool.NewStoragepoolAgent(5, ".root", "1234567890", sphc, hcMetadataType, tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/07-25-2015")
		if err != nil {
			t.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(10 * time.Second)

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

	stack = metadata.Stack{
		Name: "test_stack1",
		Services: []string{
			"service1",
			"service3",
		},
	}
	setSelfStack(stack)

	newServices := []metadata.Service{
		{
			Name: "service1",
			Containers: []string{
				"container1",
				"container2",
			},
			StackName: "test_stack1",
		},
		{
			Name: "service2",
			Containers: []string{
				"container3",
				"container4",
			},
			StackName: "test_stack2",
		},
		{
			Name: "service3",
			Containers: []string{
				"container5",
				"container6",
			},
			StackName: "test_stack1",
		},
	}
	setServices(newServices)

	newContainers := []metadata.Container{
		{
			Name:        "container1",
			ServiceName: "service1",
			StackName:   "test_stack1",
			HostUUID:    "hostUuid1",
		},
		{
			Name:        "container3",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid3",
		},
		{
			Name:        "container4",
			ServiceName: "service2",
			StackName:   "test_stack2",
			HostUUID:    "hostUuid4",
		},
	}
	setContainers(newContainers)
	time.Sleep(10 * time.Second)

	uuids = tc.getLastSync()
	if len(uuids) != 1 {
		t.Fatalf("expected 0 storagepool sync events, but received [%d] %v", len(uuids), uuids)
	}

	if uuids[0] != "hostUuid1" {
		t.Fatalf("sync event on delete does not work")
	}

}
