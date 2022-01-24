package httpclient

import (
	"net"
	"net/http"
	"time"
)

func MakeDefaultClient() *http.Client {
	return MakeNewClient(10)
}

func MakeNewClient(timeoutSeconds time.Duration) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeoutSeconds * time.Second,
				KeepAlive: timeoutSeconds * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   timeoutSeconds * time.Second,
			ResponseHeaderTimeout: timeoutSeconds * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: timeoutSeconds * time.Second,
	}
	return client
}
