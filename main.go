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
			Name:  "ip",
			Usage: "IP adresses related actions, like 'dig', 'curl', 'mtr', 'is cdn ip?' or 'ip geolocation' and so on",
			Subcommands: []cli.Command{
				{
					Name:   "is-cdn-ip",
					Usage:  "Checks whether the specified ip address is part of the Akamai edge network",
					Action: cmdCDNStatus,
				},
				{
					Name:   "geolocation",
					Usage:  "Provides the geolocation for an ip address within the Akamai network. This operation’s requests are limited to 500 per day",
					Action: cmdIPGeolocation,
				},
				{
					Name:   "dig",
					Usage:  "Run dig on a hostname to get DNS information, associating hostnames and IP addresses, from an IP address within the Akamai network not local to you",
					Action: cmdIPDig,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "hostname",
							Value: "",
							Usage: "The hostname to which to run the test",
						},
						cli.StringFlag{
							Name:  "query-type",
							Value: "A",
							Usage: "The type of DNS record, either A, AAAA, CNAME, MX, NS, PTR, or SOA. The default is A",
						},
					},
				},
				{
					Name:   "mtr",
					Usage:  "Run mtr to check connectivity between a domain and an IP address within the Akamai network",
					Action: cmdIPMtr,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "destination-domain",
							Value: "",
							Usage: "The domain name to which to test connectivity",
						},
						cli.BoolFlag{
							Name:  "resolve-dns",
							Usage: "Whether to use DNS to resolve hostnames. When disabled, output features only IP addresses",
						},
					},
				},
				{
					Name:   "curl",
					Usage:  "Run curl based on an IP address within the Akamai network. In the request object, specify a url to download and userAgent",
					Action: cmdIPCurl,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "url",
							Value: "",
							Usage: "The URL for which to gather a curl response",
						},
						cli.StringFlag{
							Name:  "user-agent",
							Value: "Chrome",
							Usage: "A header field to spoof a type of browser",
						},
					},
				},
			},
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
