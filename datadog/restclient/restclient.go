package restclient

import (
	"bytes"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	retryableClient = retryablehttp.NewClient()
	Client          HTTPClient
)

func init() {
	retryableClient.Backoff = retryablehttp.LinearJitterBackoff
	retryableClient.RetryWaitMin = 1 * time.Second
	retryableClient.RetryWaitMax = 10 * time.Second
	retryableClient.RetryMax = 2
	retryableClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryableClient.Logger = nil
	retryableClient.HTTPClient.Timeout = 30 * time.Second

	Client = retryableClient.StandardClient()
}

func Do(method string, url string, body []byte, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	request.Header = headers

	return Client.Do(request)
}
