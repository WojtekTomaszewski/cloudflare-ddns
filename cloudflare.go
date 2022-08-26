package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	cloudflareURL string = "https://api.cloudflare.com/client/v4/zones"
)

type cfClient struct {
	client *http.Client
	token  string
}

func newClient(token string) *cfClient {
	return &cfClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

func (c *cfClient) getZones() (*Zones, error) {
	req, err := http.NewRequest(http.MethodGet, cloudflareURL, nil)
	if err != nil {
		return nil, err
	}

	var zones = new(Zones)
	if err := c.request(req, &zones); err != nil {
		return nil, err
	}
	return zones, nil
}

func (c *cfClient) getDnsRecord(id, t, name string) (*Records, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/dns_records?type=%s&name=%s", cloudflareURL, id, t, name), nil)
	if err != nil {
		return nil, err
	}

	var records = new(Records)
	if err := c.request(req, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (c *cfClient) updateDnsRecord(zoneId, recordId string, record *Record) error {
	bytesReocrd, err := json.Marshal(record)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, zoneId, recordId), bytes.NewReader(bytesReocrd))
	if err != nil {
		return err
	}

	fmt.Println(req.URL)

	b, _ := json.MarshalIndent(record, "", "  ")
	fmt.Println(string(b))

	if err := c.request(req, record); err != nil {
		return err
	}
	return nil
}

func (c *cfClient) request(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(v)
}
