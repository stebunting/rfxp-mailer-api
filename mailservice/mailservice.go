package mailservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/gomail.v2"
)

type mailer struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Message       string `json:"message"`
	GreptchaToken string `json:"greptchaToken"`
	IP            string `json:"ip"`
	UserAgent     string `json:"userAgent"`
}

type recaptchaResponse struct {
	Success            bool     `json:"success"`
	ChallengeTimestamp string   `json:"challenge_ts"`
	Hostname           string   `json:"hostname"`
	Score              float32  `json:"score"`
	Action             string   `json:"action"`
	ErrorCodes         []string `json:"error-codes"`
}

type Response struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

type smtpSettings struct {
	smtpServer string
	port       int
	username   string
	password   string
	email      string
}

// getEnvVar returns a supplied environment variable for the current service
func (m *mailer) getEnvVar(name string) string {
	return os.Getenv(fmt.Sprintf("RFXP_%s", strings.ToUpper(name)))
}

// callGoogleRecaptcha makes an API call to Google and returns the response as an object.
// An error is returned if the call fails
func (m *mailer) callGoogleRecaptcha() (recaptchaResponse, error) {
	// Generate Request Body
	requestBody := url.Values{}
	requestBody.Add("secret", m.getEnvVar("RECAPTCHA_SECRET_KEY"))
	requestBody.Add("response", m.GreptchaToken)
	requestBody.Add("remoteip", m.IP)

	// Make API Call
	url := "https://www.google.com/recaptcha/api/siteverify"
	response, err := http.Post(
		url,
		"application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(requestBody.Encode())))
	if err != nil {
		return recaptchaResponse{}, err
	}
	defer response.Body.Close()

	// Parse Response Body
	var responseObject recaptchaResponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return recaptchaResponse{}, err
	}
	if err = json.Unmarshal(body, &responseObject); err != nil {
		return recaptchaResponse{}, err
	}

	return responseObject, nil
}

// // verifyGoogleRecaptcha calls Google Recaptcha Service and checks response values.
func (m *mailer) verifyGoogleRecaptcha() (bool, error) {
	var recaptchaThreshold float32 = 0.5
	response, err := m.callGoogleRecaptcha()
	if err != nil {
		return false, err
	}
	return response.Success && response.Action == "send_message" && response.Score >= recaptchaThreshold, nil
}

// Get Service Settings from environment vars
func (m *mailer) getSettings() (smtpSettings, error) {
	var s smtpSettings

	smtpPort, err := strconv.Atoi(m.getEnvVar("SMTP_PORT"))
	if err != nil {
		return s, err
	}

	s.smtpServer = m.getEnvVar("SMTP_SERVER")
	s.port = smtpPort
	s.username = m.getEnvVar("SMTP_USERNAME")
	s.password = m.getEnvVar("SMTP_PASSWORD")
	s.email = m.getEnvVar("EMAIL")

	if s.smtpServer == "" || s.port == 0 || s.username == "" || s.password == "" || s.email == "" {
		return smtpSettings{}, errors.New("Could not retrieve settings")
	}

	return s, nil
}

// sendEmail sends an email using the message args as settings
func (m *mailer) sendEmail() error {
	// Get Service Settings
	s, err := m.getSettings()
	if err != nil {
		return err
	}

	// Generate Messages
	plainTextTemplate, err := template.ParseFiles(path.Join("templates", "plaintext_email.gotmpl"))
	if err != nil {
		return err
	}
	var plainTextMsg bytes.Buffer
	plainTextTemplate.Execute(&plainTextMsg, m)

	htmlTemplate, err := template.ParseFiles(path.Join("templates", "html_email.gotmpl"))
	if err != nil {
		return err
	}
	var htmlMsg bytes.Buffer
	htmlTemplate.Execute(&htmlMsg, m)

	// // Setup Mail
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.email)
	msg.SetHeader("To", s.email)
	msg.SetHeader("Subject", "Message")
	msg.SetBody("text/plain", plainTextMsg.String())
	msg.AddAlternative("text/html", htmlMsg.String())

	// Send Mail
	d := gomail.NewDialer(s.smtpServer, s.port, s.username, s.password)
	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
