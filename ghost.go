package main

import (
	"net/url"
	"os"

	common "github.com/apiheat/akamai-cli-common"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func cmdGhostDig(c *cli.Context) error {
	return ghostDig(c)
}

func cmdGhostMtr(c *cli.Context) error {
	return ghostMtr(c)
}

func cmdGhostCurl(c *cli.Context) error {
	return ghostCurl(c)
}

func cmdGhostListLocations(c *cli.Context) error {
	return ghostListLocations(c)
}

func ghostListLocations(c *cli.Context) error {
	response, err := apiClient.ListGhostLocations()
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response.Locations))
	return nil
}

func ghostCurl(c *cli.Context) error {
	obj := common.SetStringId(c, "Please provide Ghost Location Name")

	if c.String("url") == "" {
		log.Error("Provide url, this is required parameter")
		os.Exit(4)
	}

	_, err := url.Parse(c.String("url"))
	if err != nil {
		log.Errorf("'url' is not valid URL: %s'", c.String("url"))
		os.Exit(4)
	}

	response, err := apiClient.ExecuteCurl(obj, requestFromGhost, c.String("url"), c.String("user-agent"))
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response.CurlResults))
	return nil
}

func ghostDig(c *cli.Context) error {
	obj := common.SetStringId(c, "Please provide Ghost Location Name")

	allowedQueries := []string{"A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA"}

	if c.String("hostname") == "" {
		log.Error("Provide hostname, this is required parameter")
		os.Exit(4)
	}

	u, err := url.Parse(c.String("hostname"))
	if err != nil {
		log.Errorf("'hostname' is not valid URL: %s'", c.String("hostname"))
		os.Exit(4)
	}

	if u.Scheme != "" {
		log.Errorf("Please do not provide HTTP scheme in 'hostname' : %s'", c.String("hostname"))
		os.Exit(4)
	}

	if !common.IsStringInSlice(c.String("query-type"), allowedQueries) {
		log.Error("Provided correct 'query-type': A, AAAA, CNAME, MX, NS, PTR, or SOA")
		os.Exit(5)
	}

	response, err := apiClient.ExecuteDig(obj, requestFromGhost, c.String("hostname"), c.String("query-type"))
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response.DigInfo))
	return nil
}

func ghostMtr(c *cli.Context) error {
	obj := common.SetStringId(c, "Please provide Ghost Location Name")

	if c.String("destination-domain") == "" {
		log.Error("Provide destination domain, this is required parameter")
		os.Exit(4)
	}

	u, err := url.Parse(c.String("destination-domain"))
	if err != nil {
		log.Errorf("'destination-domain' is not valid URL: %s'", c.String("destination-domain"))
		os.Exit(4)
	}

	if u.Scheme != "" {
		log.Errorf("Please do not provide HTTP scheme in 'destination-domain' : %s'", c.String("destination-domain"))
		os.Exit(4)
	}

	response, err := apiClient.ExecuteMtr(obj, requestFromGhost, c.String("destination-domain"), c.Bool("resolve-dns"))
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response.Mtr))
	return nil
}
