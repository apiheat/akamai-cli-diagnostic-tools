package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli"

	common "github.com/apiheat/akamai-cli-common"
	log "github.com/sirupsen/logrus"
)

func cmdTranslateError(c *cli.Context) error {
	return translateError(c)
}

func translateError(c *cli.Context) error {
	count := c.Int("retries")
	// Run request
	errorString := strings.Replace(common.SetStringId(c, "Please provide Error Code"), "#", "", -1)
	log.Info(fmt.Sprintf("Launch Error Translation Request for error code: %s, please note '#' is ignored", errorString))

	request, response, err := apiClient.DT.LaunchErrorTranslationRequest(errorString)
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusAccepted {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		os.Exit(2)
	}

	requestID := request.RequestID
	log.Info(fmt.Sprintf("Request for error code translation was submitted. Request ID is %s", requestID))

	// This is for making request
	// Read X-RateLimit-Remaining header, if 0 then wait for a minute with message
	// Status should be 202, if 429 - we reached limit
	if response.Response.StatusCode == http.StatusTooManyRequests {
		log.Info("Request limit per 60 seconds reached. Will wait for a minute")
		time.Sleep(61 * time.Second)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	log.Info(fmt.Sprintf("Polling error code in %d seconds", request.RetryAfter))
	time.Sleep(time.Duration(request.RetryAfter+1) * time.Second)

	s.Start()

	// Check request
	// With requestId and retryAfter data we can try to poll data
	log.Info(fmt.Sprintf("Making Translate Error request for ID: %s. Attempt 1 out of %d", requestID, c.Int("retries")))
	message, resp, err := apiClient.DT.TranslateAnError(requestID)
	count -= 2

	if resp.Response.StatusCode == http.StatusBadRequest {
		common.PrintJSON(resp.Body)
		os.Exit(3)
	}

	if err != nil || resp.Response.StatusCode != http.StatusOK {
		for {
			log.Info(fmt.Sprintf("Polling error code in %d seconds", request.RetryAfter))
			time.Sleep(time.Duration(request.RetryAfter+1) * time.Second)

			log.Info(fmt.Sprintf("Making Translate Error request for ID: %s. Attempt %d out of %d", requestID, c.Int("retries")-count, c.Int("retries")))

			count--

			message, resp, err = apiClient.DT.TranslateAnError(requestID)
			common.ErrorCheck(err)

			if resp.Response.StatusCode == http.StatusBadRequest {
				s.Stop()
				common.PrintJSON(resp.Body)
				os.Exit(3)
			}

			if resp.Response.StatusCode == http.StatusForbidden {
				s.Stop()
				common.PrintJSON(resp.Body)
				os.Exit(3)
			}

			if resp.Response.StatusCode == http.StatusOK {
				s.Stop()
				break
			}

			if count == 0 {
				s.Stop()
				log.Error("Operation took too long. Exiting...")
				os.Exit(2)
			}
		}
	}

	common.OutputJSON(message)

	return nil
}
