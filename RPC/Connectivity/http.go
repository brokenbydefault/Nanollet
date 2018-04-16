// +build !js

package Connectivity

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type H struct {
	endpoint string
}

var HTTP = H{
	endpoint: "http://[::1]:7076",
}

var client = http.Client{
	Timeout: 2 * time.Second,
}

func (c H) SendRequest(b []byte) (msg []byte, err error) {
	reader, err := c.SendRequestReader(b)
	if err != nil {
		return
	}

	return ioutil.ReadAll(reader)
}

func (c H) SendRequestReader(json []byte) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", HTTP.endpoint, bytes.NewReader(json))
	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (c H) SendRequestJSON(request interface{}, response interface{}, try ...interface{}) (err error) {
	jsn, err := json.Marshal(request)
	if err != nil {
		return
	}

	res, err := c.SendRequest(jsn)
	if err != nil {
		return
	}

	return json.Unmarshal(res, response)
}
