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
		Success: true,
	}

	t.Logf("%+v\n", expectedResponse)

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

func TestGetZonesNotOK5xx(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	expected_error := fmt.Sprintf("request failed with status code: %d", 500)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		cloudflareURL,
		httpmock.NewJsonResponderOrPanic(500, ""),
	)

	gotResponse, err := cf_client.GetZones()

	assert.Equal(t, expected_error, err.Error(), "expected error is wrong")
	assert.Nil(t, gotResponse, "response is no expected")
}

func TestGetZonesNotOK4xx(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")

	expectedResponse := &cloudflare.Zones{
		Result: []cloudflare.Zone{
			{
				ID:   "123",
				Name: "exmaple.com",
			},
		},
		Success: false,
		Errors: []cloudflare.Error{
			{
				Code:    1234,
				Message: "fail",
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		cloudflareURL,
		httpmock.NewJsonResponderOrPanic(400, expectedResponse),
	)

	gotResponse, err := cf_client.GetZones()

	// assert.Equal(t, expectedResponse.Success, gotResponse.Success, "expected success value is wrong")
	assert.Nil(t, gotResponse, "nil response expected")
	assert.EqualError(t, err, "fail")
	assert.NotNil(t, err, "error is expected")
}

func TestGetDnsRecordOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_type := "A"
	test_name := "example.com"
	test_url := fmt.Sprintf("%s/%s/dns_records?type=%s&name=%s", cloudflareURL, test_zone_id, test_type, test_name)
	expectedResponse := &cloudflare.ListRecords{
		Result: []cloudflare.Record{
			{
				Type:   test_type,
				Name:   test_name,
				ZoneID: test_zone_id,
			},
		},
		Success: true,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		test_url,
		httpmock.NewJsonResponderOrPanic(200, expectedResponse),
	)

	gotResponse, err := cf_client.GetDNSRecord(test_zone_id, test_type, test_name)

	assert.Equal(t, expectedResponse, gotResponse, "response is not correct")
	assert.Nil(t, err, "error is no expected")
}

func TestGetDnsRecordNotOK5xx(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	expected_error := fmt.Sprintf("request failed with status code: %d", 500)

	test_zone_id := "1234567"
	test_type := "A"
	test_name := "example.com"
	test_url := fmt.Sprintf("%s/%s/dns_records?type=%s&name=%s", cloudflareURL, test_zone_id, test_type, test_name)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		test_url,
		httpmock.NewJsonResponderOrPanic(500, ""),
	)

	gotResponse, err := cf_client.GetDNSRecord(test_zone_id, test_type, test_name)

	assert.Equal(t, expected_error, err.Error(), "response is not correct")
	assert.Nil(t, gotResponse, "response is no expected")
}

func TestGetDnsRecordNotOK4xx(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_type := "A"
	test_name := "example.com"
	test_url := fmt.Sprintf("%s/%s/dns_records?type=%s&name=%s", cloudflareURL, test_zone_id, test_type, test_name)
	expectedResponse := &cloudflare.ListRecords{
		Result: []cloudflare.Record{
			{
				Type:   test_type,
				Name:   test_name,
				ZoneID: test_zone_id,
			},
		},
		Success: false,
		Errors: []cloudflare.Error{
			{
				Code:    12345,
				Message: "failed",
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

	gotResponse, err := cf_client.GetDNSRecord(test_zone_id, test_type, test_name)

	assert.Nil(t, gotResponse, "nil response expected")
	assert.EqualError(t, err, "failed")
	assert.NotNil(t, err, "error is expected")
}

func TestUpdateDnsRecordOK(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_record_name := "example.com"
	test_record_id := "09876"
	test_record_type := "A"
	test_record_content := "192.168.10.10"
	test_record := cloudflare.Record{
		Name:    test_record_name,
		Type:    test_record_type,
		Content: test_record_content,
		TTL:     200,
	}
	put_record := &cloudflare.PutRecord{
		Result:  test_record,
		Success: true,
		Errors: []cloudflare.Error{
			{
				Code:    1234,
				Message: "fail",
			},
		},
	}
	test_url := fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, test_zone_id, test_record_id)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"PUT",
		test_url,
		httpmock.NewJsonResponderOrPanic(200, put_record),
	)

	err := cf_client.UpdateDNSRecord(test_zone_id, test_record_id, &test_record)

	assert.Nil(t, err, "error expected")
}

func TestUpdateDnsRecordNotOK5xx(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_record_id := "09876"
	test_record_name := "example.com"
	test_record_type := "A"
	test_record_content := "192.168.10.10"
	test_record := cloudflare.Record{
		Name:    test_record_name,
		Type:    test_record_type,
		Content: test_record_content,
		TTL:     200,
	}
	put_record := &cloudflare.PutRecord{
		Result:  test_record,
		Success: false,
		Errors: []cloudflare.Error{
			{
				Code:    1234,
				Message: "fail",
			},
		},
	}
	test_url := fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, test_zone_id, test_record_id)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"PUT",
		test_url,
		httpmock.NewJsonResponderOrPanic(500, put_record),
	)

	err := cf_client.UpdateDNSRecord(test_zone_id, test_record_id, &test_record)

	assert.NotNil(t, err, "error is expected")
	assert.EqualError(t, err, "request failed with status code: 500")

}

func TestUpdateDnsRecordNotOK4xx(t *testing.T) {
	cf_client := cloudflare.NewClient("12345")
	test_zone_id := "1234567"
	test_record_id := "09876"
	test_record_name := "example.com"
	test_record_type := "A"
	test_record_content := "192.168.10.10"
	test_record := cloudflare.Record{
		Name:    test_record_name,
		Type:    test_record_type,
		Content: test_record_content,
		TTL:     200,
	}
	put_record := &cloudflare.PutRecord{
		Result:  test_record,
		Success: false,
		Errors: []cloudflare.Error{
			{
				Code:    1234,
				Message: "fail",
			},
		},
	}
	test_url := fmt.Sprintf("%s/%s/dns_records/%s", cloudflareURL, test_zone_id, test_record_id)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"PUT",
		test_url,
		httpmock.NewJsonResponderOrPanic(400, &put_record),
	)

	err := cf_client.UpdateDNSRecord(test_zone_id, test_record_id, &test_record)

	assert.NotNil(t, err, "error is expected")
	assert.EqualError(t, err, "fail")

}
