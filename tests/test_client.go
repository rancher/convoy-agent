package tests

import (
	"fmt"

	"github.com/rancher/convoy/api"
)

type testCattleClient struct {
	lastEvents []string
	hosts      [][]string
}

func (t *testCattleClient) CreateVolume(storagepoolUuid string, vol api.VolumeResponse) error {
	t.lastEvents = append(t.lastEvents, fmt.Sprintf("CREATED_%s", vol.UUID))
	return nil
}

func (t *testCattleClient) DeleteVolume(storagepoolUuid string, vol api.VolumeResponse) error {
	t.lastEvents = append(t.lastEvents, fmt.Sprintf("DELETED_%s", vol.UUID))
	return nil
}

func (t *testCattleClient) SyncStoragePool(storagepoolUuid string, hostUuids []string) error {
	t.lastEvents = append(t.lastEvents, fmt.Sprintf("SYNC_%s", storagepoolUuid))
	t.hosts = append(t.hosts, hostUuids)
	return nil
}

func (t *testCattleClient) getLastEvent() string {
	l := len(t.lastEvents)
	if l == 0 {
		return ""
	}

	toRet := t.lastEvents[l-1]
	t.lastEvents = t.lastEvents[:l-1]

	return toRet
}

func (t *testCattleClient) getLastSync() []string {
	l := len(t.hosts)
	if l == 0 {
		return []string{}
	}

	toRet := t.hosts[l-1]
	t.hosts = t.hosts[:l-1]
	return toRet
}
