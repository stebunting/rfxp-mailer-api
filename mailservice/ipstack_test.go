package mailservice

import (
	"os"
	"testing"
)

func TestIPStackForInvalidIPs(t *testing.T) {
	tests := []string{"", "invalidip", "1.2", "34.53.132", "512.1.2.3", "1.2a.6.7"}

	for _, ip := range tests {
		m := mailer{IP: ip}

		expectedError := "invalid IP"
		err := m.getLocation()

		if err.Error() != expectedError {
			t.Error("Expected error not returned")
		}
	}
}

func TestIPStackWithInvalidKey(t *testing.T) {
	ipstackAccessKey := os.Getenv("IPSTACK_ACCESS_KEY")
	os.Setenv("IPSTACK_ACCESS_KEY", "INVALIDKEY")

	m := mailer{IP: "1.1.1.1"}
	m.getLocation()

	got := m.Location.Success
	expected := false
	if got != expected {
		t.Error("Expected error not returned")
	}

	os.Setenv("IPSTACK_ACCESS_KEY", ipstackAccessKey)
}

func TestIPStack(t *testing.T) {
	m := mailer{IP: "1.1.1.1"}
	m.getLocation()

	got := m.Location
	expected := ipStackResponse{
		CountryName: "Australia",
		City:        "Sydney",
	}
	if got.CountryName != expected.CountryName {
		t.Errorf("Unexpected country name returned from IP Stack, %s", got.CountryName)
	}
	if got.City != expected.City {
		t.Errorf("Unexpected city returned from IP Stack, %s", got.City)
	}
}
