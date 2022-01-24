package iputils

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"ipnotifier/httpclient"
	"net"
	"net/http"
	"strings"
)

func GetPublicIP() (string, error) {
	client := httpclient.MakeDefaultClient()
	res, errRes := client.Get("https://checkip.amazonaws.com")
	if errRes != nil {
		return "", errors.Wrap(errRes, "public ip http client couldn't make http request")
	}
	//req.Header.Add("User-Agent", "")

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("public ip http client responded with a non 200 code: %v", res.StatusCode))
	}

	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		return "", errors.Wrap(errBody, "public ip http client couldn't read response body")
	}

	ip := strings.TrimSuffix(string(body), "\n")
	if !IsIPValid(ip) {
		return "", errors.New("public ip http client returned an invalid ip address")
	}

	return ip, nil
}

func IsIPValid(ip string) bool {
	return net.ParseIP(ip) != nil
}
