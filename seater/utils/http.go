package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/juju/errors"
)

// http methods
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// CallAPI calls HTTP API
func CallAPI(method, url, body string, headers map[string]string) (data []byte, err error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		err = errors.Trace(err)
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("error in response: %s", string(data))
		return nil, err
	}
	return
}

// JSONAPI calls API with JSON request body and response body
func JSONAPI(method, url, body string) (j *simplejson.Json, err error) {
	data, err := CallAPI(method, url, body, map[string]string{"Content-Type": "application/json"})
	if len(data) > 0 {
		j, _ = simplejson.NewJson(data)
	} else {
		j = simplejson.New()
	}
	if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// JSONGet calls a GET JSON API
func JSONGet(url string) (j *simplejson.Json, err error) {
	return JSONAPI(MethodGet, url, "")
}

// JSONPost calls a POST JSON API
func JSONPost(url, body string) (j *simplejson.Json, err error) {
	return JSONAPI(MethodPost, url, body)
}

// JSONPut calls a PUT JSON API
func JSONPut(url, body string) (j *simplejson.Json, err error) {
	return JSONAPI(MethodPut, url, body)
}

// JSONDelete calls a DELETE JSON API
func JSONDelete(url, body string) (j *simplejson.Json, err error) {
	return JSONAPI(MethodDelete, url, body)
}

// JSONPatch calls a Patch JSON API
func JSONPatch(url, body string) (j *simplejson.Json, err error) {
	return JSONAPI(MethodPatch, url, body)
}
