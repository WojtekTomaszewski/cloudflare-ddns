package iptools

import (
	"bytes"
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
	return b.String(), nil
}

// IsIpValid does simple validation of IP address
func IsIpValid(ip string) bool {
	return net.ParseIP(ip) != nil
}
