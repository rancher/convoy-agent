package cattleevents

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"

	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"

	"github.com/rancher/convoy-agent/volume"
)

func ConnectToEventStream(conf Config) error {
	convoy, err := volume.NewConvoyClient(conf.Socket)
	log.Infof("Socket file: %v", conf.Socket)
	if err != nil {
		return err
	}

	vdh := volumeRemoveHandler{
		convoyClient: convoy,
	}
	nh := noopHandler{}
	ph := PingHandler{}

	eventHandlers := map[string]revents.EventHandler{
		"storage.volume.activate":   nh.Handler,
		"storage.volume.deactivate": nh.Handler,
		"storage.volume.remove":     vdh.Handler,
		"ping":                      ph.Handler,
	}

	router, err := revents.NewEventRouter("", 0, conf.CattleURL, conf.CattleAccessKey, conf.CattleSecretKey, nil, eventHandlers, "", conf.WorkerCount)
	if err != nil {
		return err
	}
	err = router.StartWithoutCreate(nil)
	return err
}

type volumeRemoveHandler struct {
	convoyClient *volume.ConvoyClient
}

func (h *volumeRemoveHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	data := &VSPMData{}
	err := mapstructure.Decode(event.Data, &data)
	if err != nil {
		return fmt.Errorf("Cannot parse event. Error: %v", err)
	}
	rancherVol := data.VSPM.V
	vol, err := h.convoyClient.GetVolume(rancherVol.Name)
	if err != nil {
		return fmt.Errorf("Cannot delete volume %v. Name: %v. Error: %v", rancherVol.Id, rancherVol.Name, err)
	}

	if vol == nil {
		return volumeReply(event, cli)
	}

	err = h.convoyClient.DeleteVolume(rancherVol.Name)
	if err != nil {
		return fmt.Errorf("Cannot delete volume %v. Name: %v. Error: %v", rancherVol.Id, rancherVol.Name, err)
	}

	return volumeReply(event, cli)
}

func volumeReply(event *revents.Event, cli *client.RancherClient) error {
	replyData := make(map[string]interface{})
	reply := newReply(event)
	reply.ResourceType = "volume"
	reply.ResourceId = event.ResourceId
	reply.Data = replyData
	log.Infof("Reply: %+v", reply)
	err := publishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

type noopHandler struct {
}

func (h *noopHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	log.Infof("Received and ignoring event: Name: %s, Event Id: %s, Resource Id: %s", event.Name, event.Id, event.ResourceId)
	return volumeReply(event, cli)
}

type PingHandler struct {
}

func (h *PingHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	return nil
}

func newReply(event *revents.Event) *client.Publish {
	return &client.Publish{
		Name:        event.ReplyTo,
		PreviousIds: []string{event.Id},
	}
}

func publishReply(reply *client.Publish, apiClient *client.RancherClient) error {
	_, err := apiClient.Publish.Create(reply)
	return err
}

type Config struct {
	CattleURL       string
	CattleAccessKey string
	CattleSecretKey string
	WorkerCount     int
	Socket          string
}

type VSPMData struct {
	VSPM struct {
		V struct {
			Id         int64
			Name       string
			AccountId  int64
			ExternalId string
		} `mapstructure:"volume"`
	} `mapstructure:"volumeStoragePoolMap"`
}
