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
				Name:  "storagepool-uuid",
				Usage: "set the storage pool uuid",
			},
			cli.StringFlag{
				Name:  "storagepool-healthcheck-type",
				Usage: "set the healthcheck type [file | metadata]",
				Value: "file",
			},
			cli.StringFlag{
				Name:  "storagepool-metadata-url",
				Usage: "set the metadata url",
				Value: "http://rancher-metadata/07-25-2015",
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
	storagepoolUUID := c.String("storagepool-uuid")
	if storagepoolUUID == "" {
		log.Fatalf("Required field storagepool uuid [\"storagepool-uuid\"] is not set")
	}

	storagepoolRootDir := c.GlobalString("storagepool-rootdir")

	cattleUrl := c.GlobalString("url")
	cattleAccessKey := c.GlobalString("access-key")
	cattleSecretKey := c.GlobalString("secret-key")
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	storagepoolDriver := c.GlobalString("storagepool-driver")
	storagepoolName := c.GlobalString("storagepool-name")
	if storagepoolName == "" {
		log.Fatal("required field storagepool-name has not been set")
	}

	cattleClient, err := cattle.NewCattleClient(cattleUrl, cattleAccessKey, cattleSecretKey, storagepoolDriver, storagepoolName)
	if err != nil {
		log.Fatal(err)
	}

	storagepoolAgent := NewStoragepoolAgent(healthCheckInterval, storagepoolRootDir, storagepoolUUID, healthCheckBaseDir, healthCheckType, cattleClient)

	metadataUrl := c.String("storagepool-metadata-url")

	if err := storagepoolAgent.Run(metadataUrl); err != nil {
		log.Fatal(err)
	}
}
