package iptools_test

import (
	"fmt"
	"testing"

	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/iptools"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentIPOK(t *testing.T) {
	test_ip := "192.168.10.1"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		"http://ifconfig.me",
		httpmock.NewStringResponder(200, test_ip),
	)

	ip, err := iptools.GetCurrentIP()

	assert.Equal(t, ip, test_ip, "should return ip")
	assert.Nil(t, err, "should not return error")
}

func TestGetCurrentIPNotOK(t *testing.T) {
	status_code := 400
	expected_error := fmt.Sprintf("checking current ip failed with code %d", status_code)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		"http://ifconfig.me",
		httpmock.NewStringResponder(status_code, ""),
	)

	ip, err := iptools.GetCurrentIP()

	assert.Equal(t, ip, "", "should return empty ip")
	assert.NotNil(t, err, "should return error")
	assert.Equal(t, err.Error(), expected_error, "returned error is not valid")
}

func TestIsIpValid(t *testing.T) {
	type test struct {
		in  string
		out bool
	}

	tests := []test{
		{in: "192.0.2.1", out: true},
		{in: "2001:db8::68", out: true},
		{in: "300.200.100.100", out: false},
		{in: "192.168", out: false},
	}

	for _, tc := range tests {
		out := iptools.IsIpValid(tc.in)
		if out != tc.out {
			t.Fatalf("expected %t, got %t", tc.out, out)
		}
	}
}
