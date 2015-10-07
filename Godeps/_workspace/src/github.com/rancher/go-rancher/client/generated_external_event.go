package client

const (
	EXTERNAL_EVENT_TYPE = "externalEvent"
)

type ExternalEvent struct {
	Resource

	AccountId string `json:"accountId,omitempty" yaml:"account_id,omitempty"`

	Created string `json:"created,omitempty" yaml:"created,omitempty"`

	Data map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`

	EventType string `json:"eventType,omitempty" yaml:"event_type,omitempty"`

	ExternalId string `json:"externalId,omitempty" yaml:"external_id,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	State string `json:"state,omitempty" yaml:"state,omitempty"`

	Uuid string `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type ExternalEventCollection struct {
	Collection
	Data []ExternalEvent `json:"data,omitempty"`
}

type ExternalEventClient struct {
	rancherClient *RancherClient
}

type ExternalEventOperations interface {
	List(opts *ListOpts) (*ExternalEventCollection, error)
	Create(opts *ExternalEvent) (*ExternalEvent, error)
	Update(existing *ExternalEvent, updates interface{}) (*ExternalEvent, error)
	ById(id string) (*ExternalEvent, error)
	Delete(container *ExternalEvent) error
}

func newExternalEventClient(rancherClient *RancherClient) *ExternalEventClient {
	return &ExternalEventClient{
		rancherClient: rancherClient,
	}
}

func (c *ExternalEventClient) Create(container *ExternalEvent) (*ExternalEvent, error) {
	resp := &ExternalEvent{}
	err := c.rancherClient.doCreate(EXTERNAL_EVENT_TYPE, container, resp)
	return resp, err
}

func (c *ExternalEventClient) Update(existing *ExternalEvent, updates interface{}) (*ExternalEvent, error) {
	resp := &ExternalEvent{}
	err := c.rancherClient.doUpdate(EXTERNAL_EVENT_TYPE, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ExternalEventClient) List(opts *ListOpts) (*ExternalEventCollection, error) {
	resp := &ExternalEventCollection{}
	err := c.rancherClient.doList(EXTERNAL_EVENT_TYPE, opts, resp)
	return resp, err
}

func (c *ExternalEventClient) ById(id string) (*ExternalEvent, error) {
	resp := &ExternalEvent{}
	err := c.rancherClient.doById(EXTERNAL_EVENT_TYPE, id, resp)
	return resp, err
}

func (c *ExternalEventClient) Delete(container *ExternalEvent) error {
	return c.rancherClient.doResourceDelete(EXTERNAL_EVENT_TYPE, &container.Resource)
}
