package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Cloudflare API zones endpoint
const (
	cloudflareURL string = "https://api.cloudflare.com/client/v4/zones"
)

// CFClient is representation of Cloudflare API client
type CFClient struct {
	Client *http.Client
	Token  string
}

// NewClient is CFClient constructor
func NewClient(token string) *CFClient {
	return &CFClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Token: token,
	}
}

// GetZones gets all zones DNS zones defined in the account
func (c *CFClient) GetZones() (*Zones, error) {
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

// GetDNSRecord gets specific DNS record from the zone
// id is the zone id to retriev record from
// t is record type
// name is record name, usually domain/subdomain name for which you want to change record
func (c *CFClient) GetDNSRecord(id, t, name string) (*Records, error) {
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

// UpdateDNSRecord update specific record
// zoneId is zone id where record is defined
// recordId is id of record to change
// record is payload with changes to make
func (c *CFClient) UpdateDNSRecord(zoneID, recordID string, record *Record) error {
	bytesReocrd, err := json.Marshal(record)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, zoneID, recordID), bytes.NewReader(bytesReocrd))
	if err != nil {
		return err
	}

	if err := c.request(req, record); err != nil {
		return err
	}
	return nil
}

// request does actual API call for provided req, adds rquired haeders
// res is instance of http.Request
// v represents response object for body unmarshal
func (c *CFClient) request(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(v)
}
