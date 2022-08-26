package main

import (
	"io/ioutil"
	"net/http"
)

const (
	ipURL string = "http://ifconfig.me"
)

func getCurrentIP() (string, error) {
	res, err := http.Get(ipURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
