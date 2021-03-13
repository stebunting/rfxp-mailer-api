package mailservice

import (
	"context"
	"errors"
	"fmt"
)

// returnError returns an Error status.
func returnError(err error) (Response, error) {
	return Response{
		Status:  "Error",
		Details: err.Error(),
	}, nil
}

// HandleLambdaEvent is Lambda Entry Point
func HandleLambdaEvent(ctx context.Context, m mailer) (Response, error) {
	// Verify body content
	if m.Name == "" {
		return returnError(fmt.Errorf("INVALID INPUT: '%s' is not a valid Name", m.Name))
	}
	if m.Email == "" {
		return returnError(fmt.Errorf("INVALID INPUT: '%s' is not a valid Email", m.Email))
	}
	if m.Message == "" {
		return returnError(fmt.Errorf("INVALID INPUT: '%s' is not a valid Message", m.Message))
	}
	if m.GreptchaToken == "" {
		return returnError(fmt.Errorf("INVALID INPUT: '%s' is not a valid GreptchaToken", m.GreptchaToken))
	}

	// Verify Grecaptcha Token
	validRecaptcha, err := m.verifyGoogleRecaptcha()
	if err != nil {
		return returnError(err)
	}
	if !validRecaptcha {
		return returnError(errors.New("Invalid Recaptcha"))
	}

	// Send E-Mail
	err = m.sendEmail()
	if err != nil {
		return returnError(err)
	}

	r := Response{
		Status:  "OK",
		Details: "Email sent successfully",
	}
	return r, nil
}
