package storagepool

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/rancher/convoy-agent/cattle"
)

var Commands = []cli.Command{
	{
		Name:  "storagepool",
		Usage: "Start convoy-agent as a storagepool agent",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "storagepool-healthcheck-type",
				Usage: "set the healthcheck type [file | metadata]",
				Value: "file",
			},
			cli.StringFlag{
				Name:  "storagepool-metadata-url",
				Usage: "set the metadata url",
				Value: "http://rancher-metadata/latest",
			},
		},
		Action:    storagepoolAgent,
		ShortName: "sp",
	},
}

func storagepoolAgent(c *cli.Context) {
	healthCheckInterval := c.GlobalInt("healthcheck-interval")
	healthCheckBaseDir := c.GlobalString("healthcheck-basedir")
	healthCheckType := c.String("storagepool-healthcheck-type")

	cattleUrl := c.GlobalString("url")
	cattleAccessKey := c.GlobalString("access-key")
	cattleSecretKey := c.GlobalString("secret-key")
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	storagepoolRootDir := c.GlobalString("storagepool-rootdir")
	driver := c.GlobalString("storagepool-driver")
	if driver == "" {
		log.Fatal("required field storagepool-driver has not been set")
	}

	cattleClient, err := cattle.NewCattleClient(cattleUrl, cattleAccessKey, cattleSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	storagepoolAgent := NewStoragepoolAgent(healthCheckInterval, storagepoolRootDir, driver, healthCheckBaseDir, healthCheckType, cattleClient)

	metadataUrl := c.String("storagepool-metadata-url")

	if err := storagepoolAgent.Run(metadataUrl); err != nil {
		log.Fatal(err)
	}
}
