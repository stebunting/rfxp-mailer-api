package mailservice

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	code := m.Run()
	os.Exit(code)
}

func SetEnvVars(
	email string,
	smtpServer string,
	smtpPort string,
	smtpUsername string,
	smtpPassword string,
) {
	os.Setenv("RFXP_EMAIL", email)
	os.Setenv("RFXP_SMTP_SERVER", smtpServer)
	os.Setenv("RFXP_SMTP_PORT", smtpPort)
	os.Setenv("RFXP_SMTP_USERNAME", smtpUsername)
	os.Setenv("RFXP_SMTP_PASSWORD", smtpPassword)
}

func TestInvalidInputValidation(t *testing.T) {
	HTTPClient = &http.Client{}
	var tests = []struct {
		Mailer           mailer
		ExpectedErrorMsg string
	}{
		{
			mailer{Name: "", Email: "", Message: "", GreptchaToken: ""},
			"INVALID INPUT: '' is not a valid Name",
		}, {
			mailer{Name: "John", Email: "", Message: "", GreptchaToken: ""},
			"INVALID INPUT: '' is not a valid Email",
		}, {
			mailer{Name: "Phil", Email: "phil@email.com", Message: "", GreptchaToken: ""},
			"INVALID INPUT: '' is not a valid Message",
		}, {
			mailer{Name: "Wendy", Email: "wendy@example.com", Message: "Message from Wendy", GreptchaToken: ""},
			"INVALID INPUT: '' is not a valid GreptchaToken",
		}, {
			mailer{Name: "", Email: "wendy@example.com", Message: "Message from Wendy", GreptchaToken: "FilledInToken"},
			"INVALID INPUT: '' is not a valid Name",
		}, {
			mailer{Name: "Tanaka", Email: "", Message: "Message from Tanaka", GreptchaToken: "TanakasToken"},
			"INVALID INPUT: '' is not a valid Email",
		}, {
			mailer{Name: "Honda San", Email: "hondasan@gmail.com", Message: "", GreptchaToken: "HondaSansToken"},
			"INVALID INPUT: '' is not a valid Message",
		}, {
			mailer{Name: "Valid Name", Email: "validname@example.com", Message: "Valid Message", GreptchaToken: "InvalidToken"},
			"invalid Recaptcha",
		},
	}

	context := context.Background()

	for _, test := range tests {
		response, err := HandleLambdaEvent(context, test.Mailer)
		if response.Status != "Error" {
			t.Error("Expected response status incorrect")
		}
		if response.Details != test.ExpectedErrorMsg {
			t.Error("Expected response details incorrect")
		}
		if err != nil {
			t.Error("Unexpected error returned")
		}
	}
}

func TestSmtpSettings(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"RFXP_EMAIL", "myEmail"},
		{"RFXP_SMTP_SERVER", "mySmtpServer"},
		{"RFXP_SMTP_PORT", "587"},
		{"RFXP_SMTP_USERNAME", "myUsername"},
		{"RFXP_SMTP_PASSWORD", "myPassword"},
	}

	for _, test := range tests {
		os.Setenv(test.name, test.value)
	}

	m := mailer{}
	settings, _ := m.getSettings()

	if settings.email != "myEmail" {
		t.Error("Unexpected SMTP email returned")
	}
	if settings.smtpServer != "mySmtpServer" {
		t.Error("Unexpected SMTP server returned")
	}
	if settings.port != 587 {
		t.Error("Unexpected SMTP port returned")
	}
	if settings.username != "myUsername" {
		t.Error("Unexpected SMTP username returned")
	}
	if settings.password != "myPassword" {
		t.Error("Unexpected SMTP password returned")
	}
}

func TestInvalidSmtpPortSettings(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"RFXP_SMTP_PORT", "InvalidPort"},
	}

	for _, test := range tests {
		os.Setenv(test.name, test.value)
	}

	m := mailer{}
	_, err := m.getSettings()
	if err == nil {
		t.Error("Expected an error on invalid SMTP port")
	}
}

func TestInvalidSmtpSettings(t *testing.T) {
	vars := []struct {
		name  string
		value string
	}{
		{"RFXP_EMAIL", "myEmail"},
		{"RFXP_SMTP_SERVER", "mySmtpServer"},
		{"RFXP_SMTP_PORT", "587"},
		{"RFXP_SMTP_USERNAME", "myUsername"},
		{"RFXP_SMTP_PASSWORD", "myPassword"},
	}
	for _, v := range vars {
		os.Setenv(v.name, v.value)
	}

	tests := []struct {
		name       string
		resetValue string
	}{
		{"RFXP_EMAIL", "myEmail"},
		{"RFXP_SMTP_SERVER", "mySmtpServer"},
		{"RFXP_SMTP_PORT", "587"},
		{"RFXP_SMTP_USERNAME", "myUsername"},
		{"RFXP_SMTP_PASSWORD", "myPassword"},
	}

	for _, test := range tests {
		if test.name == "RFXP_SMTP_PORT" {
			os.Setenv(test.name, "0")
		} else {
			os.Setenv(test.name, "")
		}

		m := mailer{}
		_, err := m.getSettings()
		if err.Error() != "could not retrieve settings" {
			t.Error("Unexpected error message returned")
		}
		os.Setenv(test.name, test.resetValue)
	}
}
