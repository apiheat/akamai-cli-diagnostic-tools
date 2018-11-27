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

const (
	requestFromGhost = "ghost-locations"
	requestFromIP    = "ip-addresses"
)

func main() {
	app := common.CreateNewApp(appName, "A CLI to interact with Akamai Diagnostic Tools", appVer)
	app.Flags = common.CreateFlags()

	app.Commands = []cli.Command{
		{
			Name:    "translate-request",
			Aliases: []string{"tr"},
			Usage:   "Same as 'translate-error' command, but this is not waiting for final results and you need to 'launch', 'check' and 'get' requested information",
			Subcommands: []cli.Command{
				{
					Name:      "launch",
					Usage:     "Launches a request to retrieve the data about error asynchronously. Check the poll link after the retryAfter interval or use the requestID to Check an Error Translation Request",
					UsageText: fmt.Sprintf("%s translate-error launch [command options] ERROR_CODE", appName),
					Action:    cmdLaunchTranslateErrorRequest,
				},
				{
					Name:      "check",
					Usage:     "After running 'launch' command this checks the status of an asynchronous request for data. A 200 PollResponse with a Retry-After header indicates the request is still processing. When the data is ready, a 303 response provides a Location header where you can GET the data using the 'get' command",
					UsageText: fmt.Sprintf("%s translate-error check [command options] REQUEST_ID_FROM_LAUNCH_OUTPUT", appName),
					Action:    cmdCheckTranslateErrorRequest,
				},
				{
					Name:      "get",
					Usage:     "Get information about error strings produced by edge servers when a request to retrieve content fails. The error represents an instance of a problem, and this operation gets details on what happened",
					UsageText: fmt.Sprintf("%s translate-error check [command options] REQUEST_ID_FROM_LAUNCH_OUTPUT", appName),
					Action:    cmdGetTranslateErrorRequest,
				},
			},
		},
		{
			Name:      "translate-error",
			Aliases:   []string{"t"},
			UsageText: fmt.Sprintf("%s translate-error 'Error String' --retries N", appName),
			Usage:     "Get information about error strings produced by edge servers when a request to retrieve content fails",
			Action:    cmdTranslateError,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "retries",
					Value: 50,
					Usage: "`Number` of retries to get translation request result",
				},
			},
		},
		{
			Name:  "gtm",
			Usage: "Get information about Global Traffic Management properties and gets test and target IPs for a domain and property.",
			Subcommands: []cli.Command{
				{
					Name:      "properties",
					UsageText: fmt.Sprintf("%s gtm properties", appName),
					Usage:     "List all Global Traffic Management properties (subdomains) to which you have access",
					Action:    cmdListGTMProperties,
				},
				{
					Name:      "ip-addresses",
					Usage:     "Gets test and target IPs for a domain and property. Run List GTM Properties for domain and property parameter values. PROPERTY - The Global Traffic Management property for which to collect IPs",
					UsageText: fmt.Sprintf("%s gtm ip-addresses --domain DOMAIN PROPERTY", appName),
					Action:    cmdListGTMIPAddresses,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "domain",
							Value: "",
							Usage: "The Global Traffic Management domain to which the property subdomain belongs",
						},
					},
				},
			},
		},
		{
			Name:  "ip",
			Usage: "IP addresses related actions, like 'dig', 'curl', 'mtr', 'is cdn ip?' or 'ip geolocation' and so on",
			Subcommands: []cli.Command{
				{
					Name:      "is-cdn-ip",
					UsageText: fmt.Sprintf("%s ip is-cdn-ip IP_ADDRESS", appName),
					Usage:     "Checks whether the specified ip address is part of the Akamai edge network",
					Action:    cmdCDNStatus,
				},
				{
					Name:      "geolocation",
					Usage:     "Provides the geolocation for an ip address within the Akamai network. This operation’s requests are limited to 500 per day",
					UsageText: fmt.Sprintf("%s ip geolocation IP_ADDRESS", appName),
					Action:    cmdIPGeolocation,
				},
				{
					Name:      "dig",
					Usage:     "Run dig on a hostname to get DNS information, associating hostnames and IP addresses, from an IP address within the Akamai network not local to you",
					UsageText: fmt.Sprintf("%s ip dig [command options] IP_ADDRESS", appName),
					Action:    cmdIPDig,
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
					Name:      "mtr",
					Usage:     "Run mtr to check connectivity between a domain and an IP address within the Akamai network",
					UsageText: fmt.Sprintf("%s ip mtr [command options] IP_ADDRESS", appName),
					Action:    cmdIPMtr,
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
					Name:      "curl",
					Usage:     "Run curl based on an IP address within the Akamai network. In the request object, specify a url to download and userAgent",
					UsageText: fmt.Sprintf("%s ip curl [command options] IP_ADDRESS", appName),
					Action:    cmdIPCurl,
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
			Name:  "ghost",
			Usage: "Ghost Location related actions, like 'dig', 'curl', 'mtr'",
			Subcommands: []cli.Command{
				{
					Name:      "dig",
					Usage:     "Run dig on a hostname to get DNS information, associating hostnames and IP addresses, from a location within the Akamai network not local to you. Specify location",
					UsageText: fmt.Sprintf("%s ghost dig [command options] GHOST_LOCATION", appName),
					Action:    cmdGhostDig,
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
					Name:      "locations",
					Usage:     "Lists active Akamai edge server locations from which you can run diagnostic tools",
					UsageText: fmt.Sprintf("%s ghost locations", appName),
					Action:    cmdGhostListLocations,
				},
				{
					Name:      "mtr",
					Usage:     "Run mtr to check connectivity between a domain and a location within the Akamai network not local to you. Specify location",
					UsageText: fmt.Sprintf("%s ghost mtr [command options] GHOST_LOCATION", appName),
					Action:    cmdGhostMtr,
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
					Name:      "curl",
					Usage:     "Run curl based on a location within the Akamai network. Specify location. In the request object, specify a url to download and userAgent",
					UsageText: fmt.Sprintf("%s ghost curl [command options] GHOST_LOCATION", appName),
					Action:    cmdGhostCurl,
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
					Name:      "generate",
					Usage:     "Generates a unique link to send to a user to diagnose a problem",
					UsageText: fmt.Sprintf("%s diagnostic-link generate [command options] URL", appName),
					Action:    cmdGenerateLinkRequest,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user",
							Value: "beloved-customer",
							Usage: "User name for whom you will generate link",
						},
					},
				},
				{
					Name:      "list",
					Usage:     "List users who have loaded diagnostic links over the past six months",
					UsageText: fmt.Sprintf("%s diagnostic-link list", appName),
					Action:    cmdListLinkRequests,
				},
				{
					Name:      "get",
					Usage:     "Gets details on IP addresses used for an end user’s diagnostic link test",
					UsageText: fmt.Sprintf("%s diagnostic-link get [command options] REQUEST_ID", appName),
					Action:    cmdGetLinkDetails,
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
