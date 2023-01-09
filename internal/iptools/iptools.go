package iptools

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
)

// Calling this service to get public ip
const (
	ipURL string = "http://ifconfig.me"
)

// GetCurrentIP calls ipURL for current public ip value
func GetCurrentIP() (string, error) {
	res, err := http.Get(ipURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var b bytes.Buffer
	_, err = io.Copy(&b, res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 299 {
		return "", fmt.Errorf("checking current ip failed with code %d", res.StatusCode)
	}

	return b.String(), nil
}

// IsIpValid does simple validation of IP address
func IsIPValid(ip string) bool {
	return net.ParseIP(ip) != nil
}
