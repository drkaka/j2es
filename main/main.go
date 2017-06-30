package main

import (
	"os"

	"github.com/drkaka/j2es"
	"github.com/drkaka/lg"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "j2es"
	app.Usage = "command line tool to ship logs from Journal to ElasticSearch"
	app.Version = "0.0.3"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Run in debug level.",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path for config file.",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			lg.InitLogger(true)
		} else {
			lg.InitLogger(false)
		}
		lg.L(nil).Debug("logger initialized")
		return nil
	}

	app.Action = func(c *cli.Context) error {
		return j2es.Start(c.GlobalString("config"))
	}

	app.Run(os.Args)
}
