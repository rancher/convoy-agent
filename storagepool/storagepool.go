package storagepool

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/rancher/convoy-agent/cattle"
	"github.com/rancher/convoy-agent/cattleevents"
)

var Commands = []cli.Command{
	{
		Name:  "storagepool",
		Usage: "Start convoy-agent as a storagepool agent",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "storagepool-metadata-url",
				Usage: "set the metadata url",
				Value: "http://rancher-metadata/2015-12-19",
			},
		},
		Action:    start,
		ShortName: "sp",
	},
}

func start(c *cli.Context) {
	healthCheckInterval := c.GlobalInt("healthcheck-interval")

	cattleUrl := c.GlobalString("url")
	cattleAccessKey := c.GlobalString("access-key")
	cattleSecretKey := c.GlobalString("secret-key")
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	socket := c.GlobalString("socket")

	storagepoolRootDir := c.GlobalString("storagepool-rootdir")
	driver := c.GlobalString("storagepool-driver")
	if driver == "" {
		log.Fatal("required field storagepool-driver has not been set")
	}

	cattleClient, err := cattle.NewCattleClient(cattleUrl, cattleAccessKey, cattleSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	resultChan := make(chan error)

	go func(rc chan error) {
		storagepoolAgent := NewStoragepoolAgent(healthCheckInterval, storagepoolRootDir, driver, cattleClient)
		metadataUrl := c.String("storagepool-metadata-url")
		err := storagepoolAgent.Run(metadataUrl)
		log.Errorf("Error while running storage pool agent [%v]", err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		conf := cattleevents.Config{
			CattleURL:       cattleUrl,
			CattleAccessKey: cattleAccessKey,
			CattleSecretKey: cattleSecretKey,
			WorkerCount:     10,
			Socket:          socket,
		}
		err := cattleevents.ConnectToEventStream(conf)
		log.Errorf("Cattle event listener exited with error: %s", err)
		rc <- err
	}(resultChan)

	<-resultChan
	log.Info("Exiting.")
}
