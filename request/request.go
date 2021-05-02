package request

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Request service to easy http communication
type Request struct {
	Req      *http.Request
	Res      *http.Response
	Client   *http.Client
	ResBytes []byte
	Err      error
}

// New - initiates http request with specified method, url and body
func New(method, url string, body io.Reader) *Request {
	sr := &Request{
		Client: &http.Client{Timeout: 10 * time.Second},
	}

	sr.Req, sr.Err = http.NewRequest(method, url, body)
	return sr
}

// AddHeaders - allows adding req headers
func (sr *Request) AddHeaders(key, val string) *Request {
	if sr.Err != nil {
		return sr
	}

	sr.Req.Header.Add(key, val)
	return sr
}

// Do - calls http client with defined request
func (sr *Request) Do() *Request {
	if sr.Err != nil {
		return sr
	}

	sr.Res, sr.Err = sr.Client.Do(sr.Req)
	return sr
}

// Read - reads response body into slice of bytes
func (sr *Request) Read() *Request {
	if sr.Err != nil {
		return sr
	}

	sr.ResBytes, sr.Err = ioutil.ReadAll(sr.Res.Body)
	return sr
}

// Decode - decodes response body into given data structure
func (sr *Request) Decode(data interface{}) *Request {
	if sr.Err != nil {
		return sr
	}

	sr.Err = json.NewDecoder(sr.Res.Body).Decode(data)

	return sr
}

// HasError - checks response status code and error
func (sr *Request) HasError() error {
	if sr.Err != nil {
		return sr.Err
	}

	if sr.Res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Reuest failed with status:%d", sr.Res.StatusCode)
	}

	return nil
}
