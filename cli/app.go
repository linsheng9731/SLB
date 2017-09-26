package cli

import (
	"github.com/codegangsta/cli"
	"log"
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
			Usage: "activate the silence model",
		},
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "set the filename as the configuration",
		},
	}

	app.Action = RunServer
	log.Println("run app...")
	app.Run(os.Args)

	return app
}
