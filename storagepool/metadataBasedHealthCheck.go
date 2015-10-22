package storagepool

import (
	"time"

	"github.com/rancher/go-rancher-metadata/metadata"
)

type metadataBasedHealthCheck struct {
	version     string
	prevHosts   map[string]string
	metadataUrl string
}

func (mt *metadataBasedHealthCheck) populateHostMap() (map[string]string, error) {
	m := metadata.NewClient(mt.metadataUrl)

	version, err := m.GetVersion()
	if err != nil {
		return nil, err
	}
	if version == mt.version {
		return mt.prevHosts, nil
	} else {
		mt.version = version
	}

	activeHosts := map[string]string{}
	timeStamp := time.Now().Format(time.RFC1123Z)
	stack, err := m.GetSelfStack()
	if err != nil {
		return nil, err
	}
	stackName := stack.Name

	services, err := m.GetServices()
	if err != nil {
		return nil, err
	}

	possibleServices := map[string]bool{}

	for _, service := range stack.Services {
		possibleServices[service] = true
	}

	possibleContainers := map[string]bool{}

	for _, service := range services {
		if _, ok := possibleServices[service.Name]; ok && service.StackName == stackName {
			for _, container := range service.Containers {
				possibleContainers[container] = true
			}
		}
	}

	containers, err := m.GetContainers()
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		_, okc := possibleContainers[container.Name]
		_, oks := possibleServices[container.ServiceName]
		if okc && container.StackName == stackName && oks {
			activeHosts[container.HostUUID] = timeStamp
		}
	}
	mt.prevHosts = activeHosts
	return activeHosts, nil
}

func (mt *metadataBasedHealthCheck) deleteHost(uuid string) error {
	//NoOp
	return nil
}
