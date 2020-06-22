package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/x-punch/micro-gen/cmd/new"
	"github.com/x-punch/micro-gen/cmd/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "micro-gen"
	app.Usage = "go micro service tools"
	app.Version = version.Version
	app.Commands = []*cli.Command{
		{
			Name:   "new",
			Usage:  "create empty micro service",
			Action: new.Run,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "service name", Required: true},
				&cli.StringFlag{Name: "namespace", Aliases: []string{"ns"}, Usage: "service namespace", Required: true},
				&cli.StringFlag{Name: "path", Aliases: []string{"p"}, Usage: "service generated location, default is current dir"},
			},
		},
		{
			Name:   "version",
			Usage:  "show version",
			Action: version.Run,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
