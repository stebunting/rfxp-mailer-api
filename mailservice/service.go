package mailservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/stebunting/rfxp-mailer/ipstack"
)

func init() {
	godotenv.Load()
}

// errorResponse returns an Error status.
func errorResponse(err error) (Response, error) {
	return Response{
		Status:  "Error",
		Details: err.Error(),
	}, nil
}

// HandleLambdaEvent is Lambda Entry Point
func HandleLambdaEvent(ctx context.Context, m mailer) (Response, error) {
	// Initialise Sentry
	environment := os.Getenv("ENVIRONMENT")
	err := sentry.Init(sentry.ClientOptions{
		Environment: environment,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Verify body content
	if m.Name == "" {
		return errorResponse(fmt.Errorf("INVALID INPUT: '%s' is not a valid Name", m.Name))
	}
	if m.Email == "" {
		return errorResponse(fmt.Errorf("INVALID INPUT: '%s' is not a valid Email", m.Email))
	}
	if m.Message == "" {
		return errorResponse(fmt.Errorf("INVALID INPUT: '%s' is not a valid Message", m.Message))
	}
	if m.GreptchaToken == "" {
		return errorResponse(fmt.Errorf("INVALID INPUT: '%s' is not a valid GreptchaToken", m.GreptchaToken))
	}

	// Verify Grecaptcha Token
	validRecaptcha, err := m.verifyGoogleRecaptcha()
	if err != nil {
		sentry.CaptureException(err)
		return errorResponse(err)
	}
	if !validRecaptcha {
		return errorResponse(errors.New("invalid Recaptcha"))
	}

	// Get Location Details
	m.Location, err = ipstack.GetLocation(m.IP)
	if err != nil {
		return errorResponse(err)
	}

	// Send E-Mail
	err = m.sendEmail()
	if err != nil {
		sentry.CaptureException(err)
		return errorResponse(err)
	}

	r := Response{
		Status:  "OK",
		Details: "Email sent successfully",
	}
	return r, nil
}
