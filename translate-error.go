package main

import (
	"os"
	"strings"

	common "github.com/apiheat/akamai-cli-common"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

func cmdLaunchTranslateErrorRequest(c *cli.Context) error {
	return launchErrorRequest(c)
}

func cmdCheckTranslateErrorRequest(c *cli.Context) error {
	return checkErrorRequest(c)
}

func cmdGetTranslateErrorRequest(c *cli.Context) error {
	return getErrorRequest(c)
}

func cmdTranslateError(c *cli.Context) error {
	return translateError(c)
}

func validateErrorString(err string) string {
	errorString := strings.Replace(err, "#", "", -1)
	log.Debugf("Launch Error Translation Request for error code: %s, please note '#' is ignored", errorString)

	return errorString
}

func launchErrorRequest(c *cli.Context) error {
	errorString := validateErrorString(common.SetStringId(c, "Please provide Error Code"))

	response, err := apiClient.LaunchTranslateErrorAsync(errorString)
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response))

	return nil
}

func checkErrorRequest(c *cli.Context) error {
	requestID := common.SetStringId(c, "Please provide RequestID from 'launch' command output")

	response, err := apiClient.CheckTranslateErrorAsync(requestID)
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response))

	return nil
}

func getErrorRequest(c *cli.Context) error {
	requestID := common.SetStringId(c, "Please provide RequestID from 'launch' command output")

	response, err := apiClient.RetrieveTranslateErrorAsync(requestID)
	common.ErrorCheck(err)

	common.PrintJSON(outputJSON(response))

	return nil
}

func translateError(c *cli.Context) error {
	// Run request
	errorString := validateErrorString(common.SetStringId(c, "Please provide Error Code"))

	response, err := apiClient.TranslateErrorAsync(errorString, c.Int("retries"))
	if err != nil {
		common.PrintJSON(outputJSON(err))
		os.Exit(0)
	}

	common.PrintJSON(outputJSON(response.TranslatedError))

	return nil
}
