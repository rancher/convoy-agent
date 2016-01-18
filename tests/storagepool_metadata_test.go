package tests

import (
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/rancher/convoy-agent/storagepool"
	"github.com/rancher/go-rancher-metadata/metadata"
)

var (
	hcMetadataType = "metadata"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	check.TestingT(t)
}

type MetadataTestSuite struct {
}

var _ = check.Suite(&MetadataTestSuite{})

func (s *MetadataTestSuite) SetUpSuite(c *check.C) {
	go startMetadataServer()
	time.Sleep(time.Millisecond * 200)
}

var service1 = metadata.Service{
	Name: "service1",
	Containers: []metadata.Container{
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
	},
	StackName: "test_stack1",
}

func (s *MetadataTestSuite) TestReadsHostsCorrectly(c *check.C) {
	tc := &testCattleClient{}

	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			service1,
			{
				Name: "service2",
			},
		},
	}
	setSelfStack(stack)

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(100, ".root", "1234567890", tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/mock-12-19-2015")
		if err != nil {
			c.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	uuids := tc.getLastSync()
	actual := map[string]bool{}
	for _, uuid := range uuids {
		actual[uuid] = true
	}
	c.Assert(actual, check.DeepEquals, map[string]bool{
		"hostUuid1": true,
		"hostUuid2": true,
	})
}

func (s *MetadataTestSuite) TestDetectsVersionChange(c *check.C) {
	tc := &testCattleClient{}
	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			service1,
		},
	}
	setSelfStack(stack)

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(100, ".root", "1234567890", tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/mock-12-19-2015")
		if err != nil {
			c.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	uuids := tc.getLastSync()
	actual := map[string]bool{}
	for _, uuid := range uuids {
		actual[uuid] = true
	}
	c.Assert(actual, check.DeepEquals, map[string]bool{
		"hostUuid1": true,
		"hostUuid2": true,
	})

	stack = metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			service1,
			{
				Name: "service3",
				Containers: []metadata.Container{
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
				},
				StackName: "test_stack1",
			},
		},
	}
	setSelfStack(stack)
	time.Sleep(2 * time.Second)

	uuids = tc.getLastSync()
	actual = map[string]bool{}
	for _, uuid := range uuids {
		actual[uuid] = true
	}
	c.Assert(actual, check.DeepEquals, map[string]bool{
		"hostUuid1": true,
		"hostUuid2": true,
		"hostUuid5": true,
		"hostUuid6": true,
	})
}

func (s *MetadataTestSuite) TestNoVersionChange(c *check.C) {
	tc := &testCattleClient{}
	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			service1,
			{
				Name: "service2",
			},
		},
	}
	setSelfStack(stack)

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(100, ".root", "1234567890", tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/mock-12-19-2015")
		if err != nil {
			c.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	uuids := tc.getLastSync()
	actual := map[string]bool{}
	for _, uuid := range uuids {
		actual[uuid] = true
	}
	c.Assert(actual, check.DeepEquals, map[string]bool{
		"hostUuid1": true,
		"hostUuid2": true,
	})

	stack = metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			service1,
		},
	}
	setSelfStack(stack)

	time.Sleep(10 * time.Second)

	uuids = tc.getLastSync()
	c.Assert(len(uuids), check.Equals, 0)
}

func (s *MetadataTestSuite) TestVersionChangeAndDeletion(c *check.C) {
	tc := &testCattleClient{}
	stack := metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			service1,
		},
	}
	setSelfStack(stack)

	go func() {
		spAgent := storagepool.NewStoragepoolAgent(100, ".root", "1234567890", tc)
		err := spAgent.Run("http://localhost" + metadataUrl + "/mock-12-19-2015")
		if err != nil {
			c.Fatalf("Error starting storagepool agent [%v]", err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	uuids := tc.getLastSync()
	actual := map[string]bool{}
	for _, uuid := range uuids {
		actual[uuid] = true
	}
	c.Assert(actual, check.DeepEquals, map[string]bool{
		"hostUuid1": true,
		"hostUuid2": true,
	})

	stack = metadata.Stack{
		Name: "test_stack1",
		Services: []metadata.Service{
			{
				Name: "service1",
				Containers: []metadata.Container{
					{
						Name:        "container1",
						ServiceName: "service1",
						StackName:   "test_stack1",
						HostUUID:    "hostUuid1",
					},
				},
				StackName: "test_stack1",
			},
			{
				Name:       "service3",
				Containers: []metadata.Container{},
				StackName:  "test_stack1",
			},
		},
	}
	setSelfStack(stack)
	time.Sleep(2 * time.Second)

	uuids = tc.getLastSync()
	actual = map[string]bool{}
	for _, uuid := range uuids {
		actual[uuid] = true
	}
	c.Assert(actual, check.DeepEquals, map[string]bool{
		"hostUuid1": true,
	})
}
