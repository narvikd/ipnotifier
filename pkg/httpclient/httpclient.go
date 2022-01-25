package httpclient

import (
	"net"
	"net/http"
	"time"
)

// MakeDefaultClient sets MakeNewClient() with a timeout of 10 seconds.
func MakeDefaultClient() *http.Client {
	return MakeNewClient(10)
}

// MakeNewClient returns a pointer to a new client which has all timeouts set to timeoutSeconds,
// except Transport.ExpectContinueTimeout which is set to "1 * time.Second" .
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
