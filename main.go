package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/rancher/convoy-agent/storagepool"
	"github.com/rancher/convoy-agent/volume"
	"github.com/rancher/kubernetes-agent/healthcheck"
)

var (
	VERSION = "0.10.0-dev"
	port    = 10241
)

func main() {
	app := cli.NewApp()
	app.Name = "convoy-agent"
	app.Version = VERSION
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
			Value: 5000,
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
			Usage: "set the storage pool driver.",
		},
		cli.StringFlag{
			Name:  "socket, s",
			Value: "/var/run/convoy/convoy.sock",
			Usage: "specify unix domain socket for communicating with convoy server",
		},
	}

	commands := append(volume.Commands, storagepool.Commands...)
	app.Commands = commands

	go func() {
		err := healthcheck.StartHealthCheck(port)
		log.Fatalf("Error while running healthcheck [%v]", err)
	}()
	app.EnableBashCompletion = true
	app.Run(os.Args)
}
