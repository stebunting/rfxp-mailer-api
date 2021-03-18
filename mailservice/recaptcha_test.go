package mailservice

import (
	"testing"

	"github.com/stebunting/rfxp-mailer/mocks"
)

func TestRecaptcha(t *testing.T) {
	m := mailer{}
	HTTPClient = &mocks.MockHTTPClient{
		Resp: recaptchaResponse{
			Success: true,
		},
	}
	got, _ := m.callGoogleRecaptcha()
	if !got.Success {
		t.Error("Unexpected status")
	}
}

func TestVerifyRecaptcha(t *testing.T) {
	tests := []struct {
		RecaptchaResponse recaptchaResponse
		ExpectedResponse  bool
		ErrorReturned     bool
	}{
		{recaptchaResponse{Success: true, Score: 0.9, Action: "send_message"}, true, false},
		{recaptchaResponse{Success: true, Score: 0.5, Action: "send_message"}, true, false},
		{recaptchaResponse{Success: true, Score: 0.4, Action: "send_message"}, false, false},
		{recaptchaResponse{Success: false, Score: 0.9, Action: "send_message"}, false, false},
		{recaptchaResponse{Success: true, Score: 0.9, Action: "invalid_message"}, false, false},
		{recaptchaResponse{Success: true, Score: 0.9, Action: "send_message", ErrorCodes: []string{"101", "102"}}, false, true},
	}

	for _, test := range tests {
		m := mailer{}
		HTTPClient = &mocks.MockHTTPClient{
			Error: len(test.RecaptchaResponse.ErrorCodes) > 0,
			Resp:  test.RecaptchaResponse,
		}
		got, err := m.verifyGoogleRecaptcha()

		if got != test.ExpectedResponse {
			t.Error("Expected response incorrect")
		}
		if test.ErrorReturned == (err == nil) {
			t.Error("Unexpected error returned")
		}
	}
}
