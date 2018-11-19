package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	common "github.com/apiheat/akamai-cli-common"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func cmdCDNStatus(c *cli.Context) error {
	return isCDNIP(c)
}

func isCDNIP(c *cli.Context) error {
	ip := common.SetStringId(c, "Please provide IP")

	if !isIPv4(ip) {
		log.Info("Provided IP address is not valid IPv4 address:", ip)
		os.Exit(3)
	}

	request, response, err := apiClient.DT.CDNStatus(ip)
	common.ErrorCheck(err)

	if response.Response.StatusCode != 200 {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		os.Exit(2)
	}

	common.OutputJSON(request)

	return nil
}

func isIPv4(host string) bool {
	parts := strings.Split(host, ".")

	if len(parts) < 4 {
		return false
	}

	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}
