package ipstack

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stebunting/rfxp-mailer/mocks"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	code := m.Run()
	os.Exit(code)
}

func TestIPStackForInvalidIPs(t *testing.T) {
	tests := []string{"", "invalidip", "1.2", "34.53.132", "512.1.2.3", "1.2a.6.7"}

	for _, ip := range tests {
		expectedError := "invalid IP"
		_, err := GetLocation(ip)

		if err.Error() != expectedError {
			t.Error("Expected error not returned")
		}
	}
}

func TestIPStackWithInvalidKey(t *testing.T) {
	ipstackAccessKey := os.Getenv("IPSTACK_ACCESS_KEY")
	os.Setenv("IPSTACK_ACCESS_KEY", "INVALIDKEY")

	_, err := GetLocation("1.1.1.1")
	expectedErr := "IP Stack call failed"
	if err.Error() != expectedErr {
		t.Error("Expected error not returned")
	}

	os.Setenv("IPSTACK_ACCESS_KEY", ipstackAccessKey)
}

func TestIPStackWthGetError(t *testing.T) {
	httpClientBackup := HTTPClient

	HTTPClient = &mocks.MockHTTPClient{Error: true}
	_, err := GetLocation("1.1.1.1")
	if err == nil {
		t.Error("Expected error not returned")
	}

	HTTPClient = httpClientBackup
}

func TestIPStackRealCall(t *testing.T) {
	got, err := GetLocation("14.90.14.201")
	if err != nil {
		t.Error(err)
		t.Error("Unexpected error")
	}

	expected := Location{
		CountryName: "South Korea",
		City:        "Seongnam-si",
	}
	if got.CountryName != expected.CountryName {
		t.Errorf("Unexpected country name returned from IP Stack, %s", got.CountryName)
	}
	if got.City != expected.City {
		t.Errorf("Unexpected city returned from IP Stack, %s", got.City)
	}
}
