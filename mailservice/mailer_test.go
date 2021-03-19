package mailservice

import (
	"testing"

	"github.com/stebunting/rfxp-mailer/mocks"
)

func TestMailer(t *testing.T) {
	SetEnvVars("myemail@example.com", "smtp.example.com", "587", "myUsername", "myPassword")
	HTTPClient = &mocks.MockHTTPClient{
		Resp: recaptchaResponse{
			Success: true,
			Score:   0.8,
			Action:  "send_message",
		},
	}
	m := mailer{
		Name:          "Test User",
		Email:         "test@example.com",
		Message:       "This is a test message.",
		GreptchaToken: "INVALIDTOKEN",
		IP:            "141.65.181.161",
		UserAgent:     "Go Test",
	}
	err := m.sendEmail()
	if err == nil {
		t.Error(err)
	}
}
