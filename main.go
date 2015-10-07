package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/rancher/convoy-agent/storagepool"
	"github.com/rancher/convoy-agent/volume"
)

var (
	GITCOMMIT = "HEAD"
)

func main() {
	app := cli.NewApp()
	app.Name = "convoy-agent"
	app.Version = GITCOMMIT
	app.Author = "Rancher Labs"
	app.Usage = "An agent that acts as an interface between rancher storage and cattle server"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug logging level",
		},
		cli.StringFlag{
			Name:   "url",
			Usage:  "The URL endpoint to communicate with cattle server",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  "The access key required to authenticate with cattle server",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "The secret key required to authenticate with cattle server",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:  "healthcheck-interval",
			Value: 5,
			Usage: "set the frequency of performing healthchecks",
		},
		cli.StringFlag{
			Name:  "healthcheck-basedir",
			Value: ".healthcheck",
			Usage: "set the directory to write the healthcheck files into",
		},
		cli.StringFlag{
			Name:  "storagepool-rootdir",
			Usage: "set the storage pool rootdir",
			Value: ".root",
		},
		cli.StringFlag{
			Name:  "storagepool-driver",
			Usage: "set the storage pool driver",
			Value: "convoy",
		},
		cli.StringFlag{
			Name:  "storagepool-name",
			Usage: "set the storage pool name",
		},
	}

	commands := append(volume.Commands, storagepool.Commands...)
	app.Commands = commands

	app.EnableBashCompletion = true
	app.Run(os.Args)
}
