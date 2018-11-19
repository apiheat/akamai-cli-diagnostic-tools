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
					Usage: "Number of retries to get translation request result",
				},
			},
		},
		{
			Name:    "is-akamai-cdn-ip",
			Aliases: []string{"i"},
			Usage:   "Checks whether the specified ip address is part of the Akamai edge network",
			Action:  cmdCDNStatus,
		},
		{
			Name:  "diagnostic-link",
			Usage: "Generate/List/Get a unique link to send to a user to diagnose a problem",
			Subcommands: []cli.Command{
				{
					Name:   "generate",
					Usage:  "Generates a unique link to send to a user to diagnose a problem",
					Action: cmdGenerateLinkRequest,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user",
							Value: "beloved-customer",
							Usage: "User name for whom you will generate link",
						},
					},
				},
				{
					Name:   "list",
					Usage:  "List users who have loaded diagnostic links over the past six months",
					Action: cmdListLinkRequests,
				},
				{
					Name:   "get",
					Usage:  "Gets details on IP addresses used for an end user’s diagnostic link test",
					Action: cmdGetLinkDetails,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "retries",
							Value: 10,
							Usage: "Number of retries to get ip address request result",
						},
					},
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
