package main

import (
	"fmt"
	"os"
	"sort"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	apiClient       *edgegrid.Client
	appName, appVer string
)

func main() {
	app := common.CreateNewApp(appName, "A CLI to interact with Akamai Diagnostic Tools", appVer)
	app.Flags = common.CreateFlags()

	app.Commands = []cli.Command{
		{
			Name:    "translate-error",
			Aliases: []string{"t"},
			Usage:   "Get information about error strings produced by edge servers when a request to retrieve content fails",
			Action:  cmdTranslateError,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "retries",
					Value: 50,
					Usage: "Number of retries to get translation request results",
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Before = func(c *cli.Context) error {
		var err error

		apiClient, err = common.EdgeClientInit(c.GlobalString("config"), c.GlobalString("section"), c.GlobalString("debug"))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
