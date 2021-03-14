package mailservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	RecaptchaResponse recaptchaResponse
}

func (m *mockHTTPClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	if len(m.RecaptchaResponse.ErrorCodes) > 0 {
		err = fmt.Errorf("%v", m.RecaptchaResponse.ErrorCodes)
	}
	replyString, _ := json.Marshal(m.RecaptchaResponse)

	resp = &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(string(replyString))),
	}
	return
}

func TestRecaptcha(t *testing.T) {
	m := mailer{}
	HTTPClient = &mockHTTPClient{
		recaptchaResponse{
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
		HTTPClient = &mockHTTPClient{test.RecaptchaResponse}
		got, err := m.verifyGoogleRecaptcha()

		if got != test.ExpectedResponse {
			t.Error("Expected response incorrect")
		}
		if test.ErrorReturned == (err == nil) {
			t.Error("Unexpected error returned")
		}
	}
}
