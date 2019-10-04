package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

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

	response, err := apiClient.GenerateDiagnosticLink(c.String("user"), testURL)
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response))

	return nil
}

func listLinkRequests(c *cli.Context) error {
	response, err := apiClient.ListDiagnosticLinkRequests()
	common.ErrorCheck(err)

	common.OutputJSON(response.EndUserIPRequests)
	return nil
}

func getLinkRequest(c *cli.Context) error {
	requestID := common.SetStringId(c, "Please provide valid Request ID")

	response, err := apiClient.RetrieveDiagnosticLinkRequest(requestID)
	common.ErrorCheck(err)

	common.OutputJSON(response.EndUserIPDetails)

	return nil
}

func unescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func outputJSON(input interface{}) string {
	b, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err)
	}

	jsonRawUnescaped, _ := unescapeUnicodeCharactersInJSON(json.RawMessage(b))

	return string(jsonRawUnescaped)
}
