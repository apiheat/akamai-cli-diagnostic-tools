package main

import (
	"fmt"
	"net/http"
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
	request, response, err := apiClient.DT.ListGTMProperties()
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		common.PrintJSON(response.Body)
		os.Exit(2)
	}

	common.OutputJSON(request.GtmProperties)

	return nil
}

func listGTMIPs(c *cli.Context) error {
	property := common.SetStringId(c, "Please provide PROPERTY. The Global Traffic Management property for which to collect IPs")

	if c.String("domain") == "" {
		log.Error("Provide domain, this is required parameter. The Global Traffic Management domain to which the property subdomain belongs")
		os.Exit(4)
	}

	request, response, err := apiClient.DT.ListGTMPropertyIPs(property, c.String("domain"))
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		common.PrintJSON(response.Body)
		os.Exit(2)
	}

	common.OutputJSON(request.GtmPropertyIps)

	return nil
}
