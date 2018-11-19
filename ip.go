package main

import (
	"fmt"
	"net/http"
	"net/url"
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

func cmdIPGeolocation(c *cli.Context) error {
	return ipGeolocation(c)
}

func cmdIPDig(c *cli.Context) error {
	return ipDig(c)
}

func cmdIPMtr(c *cli.Context) error {
	return ipMtr(c)
}

func cmdIPCurl(c *cli.Context) error {
	return ipCurl(c)
}

func ipCurl(c *cli.Context) error {
	obj := common.SetStringId(c, "Please provide IP")

	if !isIPv4(obj) {
		log.Error("Provided IP address is not valid IPv4 address:", obj)
		os.Exit(3)
	}

	if c.String("url") == "" {
		log.Error("Provide url, this is required parameter")
		os.Exit(4)
	}

	_, err := url.Parse(c.String("url"))
	if err != nil {
		log.Errorf("'url' is not valid URL: %s'", c.String("url"))
		os.Exit(4)
	}

	_, response, err := apiClient.DT.Curl(obj, requestFromIP, c.String("url"), c.String("user-agent"))
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		os.Exit(2)
	}

	common.PrintJSON(response.Body)

	return nil
}

func ipMtr(c *cli.Context) error {
	ip := common.SetStringId(c, "Please provide IP")

	if !isIPv4(ip) {
		log.Error("Provided IP address is not valid IPv4 address:", ip)
		os.Exit(3)
	}

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

	_, response, err := apiClient.DT.IPMtr(ip, c.String("destination-domain"), c.Bool("resolve-dns"))
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		os.Exit(2)
	}

	common.PrintJSON(response.Body)

	return nil
}

func ipDig(c *cli.Context) error {
	obj := common.SetStringId(c, "Please provide IP")

	allowedQueries := []string{"A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA"}

	if !isIPv4(obj) {
		log.Error("Provided IP address is not valid IPv4 address:", obj)
		os.Exit(3)
	}

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

	_, response, err := apiClient.DT.Dig(obj, requestFromIP, c.String("hostname"), c.String("query-type"))
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		os.Exit(2)
	}

	common.PrintJSON(response.Body)

	return nil
}

func ipGeolocation(c *cli.Context) error {
	ip := common.SetStringId(c, "Please provide IP")

	if !isIPv4(ip) {
		log.Info("Provided IP address is not valid IPv4 address:", ip)
		os.Exit(3)
	}

	request, response, err := apiClient.DT.IPGeolocation(ip)
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Something went wrong, re-run in debug mode. Response code: %d", response.Response.StatusCode))
		os.Exit(2)
	}

	common.OutputJSON(request.GeoLocation)

	return nil
}

func isCDNIP(c *cli.Context) error {
	ip := common.SetStringId(c, "Please provide IP")

	if !isIPv4(ip) {
		log.Info("Provided IP address is not valid IPv4 address:", ip)
		os.Exit(3)
	}

	request, response, err := apiClient.DT.CDNStatus(ip)
	common.ErrorCheck(err)

	if response.Response.StatusCode != http.StatusOK {
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
