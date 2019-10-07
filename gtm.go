package main

import (
	"os"

	common "github.com/apiheat/akamai-cli-common"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func cmdListGTMProperties(c *cli.Context) error {
	return listGTM(c)
}

func cmdListGTMIPAddresses(c *cli.Context) error {
	return listGTMIPs(c)
}

func listGTM(c *cli.Context) error {
	response, err := apiClient.ListGTMProperties()
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response.GtmProperties))
	return nil
}

func listGTMIPs(c *cli.Context) error {
	property := common.SetStringId(c, "Please provide PROPERTY. The Global Traffic Management property for which to collect IPs")

	if c.String("domain") == "" {
		log.Error("Provide domain, this is required parameter. The Global Traffic Management domain to which the property subdomain belongs")
		os.Exit(4)
	}

	response, err := apiClient.ListGTMPropertyIPs(property, c.String("domain"))
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response.GtmPropertyIps))

	return nil
}
