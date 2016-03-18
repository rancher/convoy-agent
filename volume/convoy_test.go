package volume

import (
	"testing"

	"gopkg.in/check.v1"
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
}

var _ = check.Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *check.C) {

}

func (s *TestSuite) TestVolumeDelete(c *check.C) {
	name := "foo"
	convoyClient, err := NewConvoyClient(testSock)
	if err != nil {
		c.Fatal(err)
	}
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
	err = convoyClient.DeleteVolume(name)
	if err != nil {
		c.Fatal(err)
	}
	vols, err = convoyClient.GetCurrVolumes()
	if err != nil {
		c.Fatal(err)
	}
	found = false
	for _, vol := range vols {
		if vol.Name == name {
			found = true
			break
		}
	}
	if found {
		c.Fatalf("Volume %v was not deleted.", name)
	}
}
