package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
)

// Request struct
type Request struct {
	client *httpclient.Client
}

// HTTPRequest interface
type HTTPRequest interface {
	Do(context context.Context, method, url string, body io.Reader, target interface{}, headers map[string]string) (string, error)
}

// NewRequest function
// Request's Constructor
// Returns : *Request
func NewRequest(retries int, sleepBetweenRetry time.Duration) HTTPRequest {
	// define a maximum jitter interval
	maximumJitterInterval := 5 * time.Millisecond

	// create a backoff
	backoff := heimdall.NewConstantBackoff(sleepBetweenRetry, maximumJitterInterval)

	// create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	// set http timeout
	timeout := 2000 * time.Millisecond

	// set http client
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(retries),
	)

	return &Request{
		client: client,
	}
}

// Do function, for http client call
func (request *Request) Do(context context.Context, method, url string, body io.Reader, target interface{}, headers map[string]string) (string, error) {
	var (
		respBody   []byte
		respStatus string
	)

	// set request http
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return respStatus, err
	}

	// iterate optional data of headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// client request
	r, err := request.client.Do(req)
	if err != nil {
		return respStatus, err
	}

	// close response body
	defer r.Body.Close()

	// set resp status
	respStatus = r.Status

	// read to response body with io util
	respBody, err = ioutil.ReadAll(r.Body)

	// unmarshal to our target
	err = json.Unmarshal(respBody, target)
	if err != nil {
		return respStatus, err
	}

	return respStatus, nil
}

