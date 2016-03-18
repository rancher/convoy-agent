package cattleevents

import (
	"testing"

	"gopkg.in/check.v1"

	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"

	"github.com/rancher/convoy-agent/volume"
)

const testSock string = "/var/run/convoy/convoy.sock"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	check.TestingT(t)
}

// These tests presume an instance of conovy running and exposed on socket
// /var/run/test/convoy-test.sock and registered with docker convoy-test
// The test scripts set this up.
type TestSuite struct {
	publishChan chan client.Publish
	mockRClient *client.RancherClient
}

var _ = check.Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *check.C) {
	s.publishChan = make(chan client.Publish, 10)

	mock := &MockPublishOperations{
		publishChan: s.publishChan,
	}
	s.mockRClient = &client.RancherClient{
		Publish: mock,
	}
}

func (s *TestSuite) TestVolumeDelete(c *check.C) {
	convoyClient, err := volume.NewConvoyClient(testSock)
	if err != nil {
		c.Fatal(err)
	}

	handler := volumeRemoveHandler{
		convoyClient: convoyClient,
	}

	name := "handlertest"
	err = convoyClient.CreateVolume(name)
	if err != nil {
		c.Fatal(err)
	}
	vols, err := convoyClient.GetCurrVolumes()
	if err != nil {
		c.Fatal(err)
	}
	found := false
	for _, vol := range vols {
		if vol.Name == name {
			found = true
			break
		}
	}
	if !found {
		c.Fatalf("Volume %v was not created.", name)
	}

	event := &revents.Event{
		ReplyTo: "event-1",
		Id:      "event-id-1",
		Data: map[string]interface{}{
			"volumeStoragePoolMap": &map[string]interface{}{
				"volume": &map[string]interface{}{
					"name":      name,
					"id":        1,
					"accountId": 1,
				},
			},
		},
	}

	err = handler.Handler(event, s.mockRClient)
	if err != nil {
		c.Fatal(err)
	}
	pub := <-s.publishChan
	c.Assert(pub.Name, check.Equals, "event-1")
	c.Assert(pub.PreviousIds, check.DeepEquals, []string{"event-id-1"})
	c.Assert(len(pub.Data), check.Equals, 0)

	// Assert that the event running a second time does not fail.
	err = handler.Handler(event, s.mockRClient)
	if err != nil {
		c.Fatal(err)
	}
	pub = <-s.publishChan
	c.Assert(pub.Name, check.Equals, "event-1")
	c.Assert(pub.PreviousIds, check.DeepEquals, []string{"event-id-1"})
	c.Assert(len(pub.Data), check.Equals, 0)
}

type MockPublishOperations struct {
	client.PublishClient
	publishChan chan<- client.Publish
}

func (m *MockPublishOperations) Create(publish *client.Publish) (*client.Publish, error) {
	m.publishChan <- *publish
	return nil, nil
}
