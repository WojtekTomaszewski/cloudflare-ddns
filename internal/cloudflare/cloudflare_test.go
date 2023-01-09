package cloudflare_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/cloudflare"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	cloudflareURL string = "https://api.cloudflare.com/client/v4/zones"
)

func TestNewClient(t *testing.T) {
	test_token := "12345"

	expected_client := &cloudflare.CFClient{
		Token: test_token,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	got_client := cloudflare.NewClient(test_token)

	assert.Equal(t, *expected_client, *got_client, "returned client is not the same as expected")
	assert.Equal(t, test_token, got_client.Token, "returned client token is wrong")

}

func TestGetZonesOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	expectedResponse := &cloudflare.Zones{
		Result: []cloudflare.Zone{
			{
				ID:   "123",
				Name: "exmaple.com",
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		cloudflareURL,
		httpmock.NewJsonResponderOrPanic(200, expectedResponse),
	)

	gotResponse, err := cf_client.GetZones()

	assert.Equal(t, expectedResponse, gotResponse, "expected zone list is wrong")
	assert.Nil(t, err, "error is no expected")
}

func TestGetZonesNotOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	expected_error := fmt.Sprintf("request failed with status code: %d", 400)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		cloudflareURL,
		httpmock.NewJsonResponderOrPanic(400, ""),
	)

	gotResponse, err := cf_client.GetZones()

	assert.Equal(t, expected_error, err.Error(), "expected error is wrong")
	assert.Nil(t, gotResponse, "response is no expected")
}

func TestGetDnsRecordOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_type := "A"
	test_name := "example.com"
	test_url := fmt.Sprintf("%s/%s/dns_records?type=%s&name=%s", cloudflareURL, test_zone_id, test_type, test_name)
	expectedResponse := &cloudflare.Records{
		Result: []cloudflare.Record{
			{
				Type:   test_type,
				Name:   test_name,
				ZoneID: test_zone_id,
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		test_url,
		httpmock.NewJsonResponderOrPanic(200, expectedResponse),
	)

	gotResponse, err := cf_client.GetDnsRecord(test_zone_id, test_type, test_name)

	assert.Equal(t, expectedResponse, gotResponse, "response is not correct")
	assert.Nil(t, err, "error is no expected")
}

func TestGetDnsRecordNotOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	expected_error := fmt.Sprintf("request failed with status code: %d", 403)

	test_zone_id := "1234567"
	test_type := "A"
	test_name := "example.com"
	test_url := fmt.Sprintf("%s/%s/dns_records?type=%s&name=%s", cloudflareURL, test_zone_id, test_type, test_name)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		test_url,
		httpmock.NewJsonResponderOrPanic(403, ""),
	)

	gotResponse, err := cf_client.GetDnsRecord(test_zone_id, test_type, test_name)

	assert.Equal(t, expected_error, err.Error(), "response is not correct")
	assert.Nil(t, gotResponse, "response is no expected")
}

func TestUpdateDnsRecordOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_record_name := "example.com"
	test_record_id := "09876"
	test_record_type := "A"
	test_record_content := "192.168.10.10"
	test_record := &cloudflare.Record{
		Name:    test_record_name,
		Type:    test_record_type,
		Content: test_record_content,
		TTL:     200,
	}
	test_url := fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, test_zone_id, test_record_id)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"PUT",
		test_url,
		httpmock.NewJsonResponderOrPanic(200, test_record),
	)

	err := cf_client.UpdateDnsRecord(test_zone_id, test_record_id, test_record)

	assert.Nil(t, err, "no error expected")

}

func TestUpdateDnsRecordNotOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_record_id := "09876"
	test_record := &cloudflare.Record{}
	test_url := fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, test_zone_id, test_record_id)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"PUT",
		test_url,
		httpmock.NewJsonResponderOrPanic(400, ""),
	)

	err := cf_client.UpdateDnsRecord(test_zone_id, test_record_id, test_record)

	assert.NotNil(t, err, "error is expected")

}
