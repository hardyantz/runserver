package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func httpRequest(method, url string, body io.Reader, headers map[string]string, target interface{}) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	// iterate optional data of headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)

	// unmarshal to our target
	err = json.Unmarshal(respBody, target)
	if err != nil {
		return err
	}

	return nil
}

