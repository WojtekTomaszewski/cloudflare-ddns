package publicip

import (
	"io/ioutil"
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
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// IsIpValid does simple validation of IP address
func IsIpValid(ip string) bool {
	return net.ParseIP(ip) != nil
}
