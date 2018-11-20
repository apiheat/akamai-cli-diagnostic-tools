package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	common "github.com/apiheat/akamai-cli-common"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func cmdGenerateLinkRequest(c *cli.Context) error {
	return generateLink(c)
}

func cmdListLinkRequests(c *cli.Context) error {
	return listLinkRequests(c)
}

func cmdGetLinkDetails(c *cli.Context) error {
	return getLinkRequest(c)
}

func generateLink(c *cli.Context) error {
	testURL := common.SetStringId(c, "Please provide URL you want to simulate the user loading")

	_, err := url.Parse(testURL)
	if err != nil {
		log.Error(fmt.Sprintf("URL you want to simulate the user loading is not valid URL '%s'", testURL))
		os.Exit(3)
	}

	_, response, err := apiClient.DT.GenerateDiagnosticLink(c.String("user"), testURL)
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusCreated {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		common.PrintJSON(response.Body)
		os.Exit(2)
	}

	common.PrintJSON(response.Body)

	return nil
}

func listLinkRequests(c *cli.Context) error {
	request, response, err := apiClient.DT.ListDiagnosticLinkRequests()
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		common.PrintJSON(response.Body)
		os.Exit(2)
	}

	common.OutputJSON(request.EndUserIPRequests)
	return nil
}

func getLinkRequest(c *cli.Context) error {
	requestID := common.SetStringId(c, "Please provide valid Request ID")

	request, response, err := apiClient.DT.GetDiagnosticLinkRequest(requestID)
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		common.PrintJSON(response.Body)
		os.Exit(2)
	}

	common.OutputJSON(request.EndUserIPDetails)

	return nil
}
