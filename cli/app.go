package cli

import (
	"github.com/codegangsta/cli"
	"os"
)

const APP_NAME = "SLB (github.com/linsheng9731/SLB)"
const APP_USAGE = "slb"
const VERSION_MAJOR = "0"
const VERSION_MINOR = "1"
const VERSION_BUILD = "0"

func CreateAPP() *cli.App {
	app := cli.NewApp()

	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_BUILD

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "silence, s",
			Usage: "drop verbose log information ",
		},
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "set the filename as the configuration",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "reload",
			Usage:  "reload configure without downtime",
			Action: HotReload,
		},
	}
	app.Action = RunServer
	app.Run(os.Args)

	return app
}
