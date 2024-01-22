package proxy

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

type ProxyData struct {
	TestUrl string
	TimeOut int
	Data    interface{}
}

func CheckProxyValidity(proxyAddress string, timeout time.Duration, testURL string) bool {

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse("http://" + proxyAddress)
			},
			Dial: (&net.Dialer{
				Timeout: timeout * time.Second,
			}).Dial,
		},
		Timeout: timeout * time.Second,
	}

	resp, err := httpClient.Get(testURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
