package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	methodError = "HTTP Method is not suppost!"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
)

func httpRequest(u string, method string, data map[string]string) (body []byte, err error) {
	var (
		req       *http.Request
		resp      *http.Response
		targetURL *url.URL
	)
	// Data or Params
	d := url.Values{}
	for k, v := range data {
		d.Set(k, v)
	}
	if method == "GET" {
		if targetURL, err = url.Parse(u); err != nil {
			return nil, err
		}
		targetURL.RawQuery = d.Encode()
		if req, err = http.NewRequest(method, targetURL.String(), nil); err != nil {
			return nil, err
		}
	} else if method == "POST" {
		if req, err = http.NewRequest(method, u, strings.NewReader(d.Encode())); err != nil {
			return nil, err
		}
	} else if method == "JSON" {
		var jsonByte []byte
		if _, ok := data["json"]; ok {
			jsonByte = []byte(data["json"])
		} else {
			if jsonByte, err = json.Marshal(data); err != nil {
				return nil, err
			}
		}
		if req, err = http.NewRequest("POST", u, bytes.NewBuffer(jsonByte)); err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		return nil, errors.New(methodError)
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

func httpGet(u string, params map[string]string) (body []byte, err error) {
	return httpRequest(u, "GET", params)
}

func httpPost(u string, data map[string]string) (body []byte, err error) {
	return httpRequest(u, "POST", data)
}

func httpJSON(u string, data map[string]string) (body []byte, err error) {
	return httpRequest(u, "JSON", data)
}
