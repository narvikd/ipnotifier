package iputils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ipnotifier/pkg/errorsutils"
	"ipnotifier/pkg/httpclient"
	"net"
	"net/http"
	"strings"
)

// GetPublicIP returns the public IP from "https://checkip.amazonaws.com".
func GetPublicIP() (string, error) {
	const endpoint = "https://checkip.amazonaws.com"
	client := httpclient.MakeDefaultClient()
	res, errRes := client.Get(endpoint)
	if errRes != nil {
		return "", errorsutils.Wrap(errRes, "public ip http client couldn't make http request")
	}
	//req.Header.Add("User-Agent", "")

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("public ip http client responded with a non 200 code: %v", res.StatusCode)
		return "", errors.New(errStr)
	}

	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		return "", errorsutils.Wrap(errBody, "public ip http client couldn't read response body")
	}

	ip := strings.TrimSuffix(string(body), "\n")
	if !IsIPValid(ip) {
		return "", errors.New("public ip http client returned an invalid ip address")
	}

	return ip, nil
}

// IsIPValid returns a boolean indicating if the IP is valid.
func IsIPValid(ip string) bool {
	return net.ParseIP(ip) != nil
}
