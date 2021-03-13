package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stebunting/rfxp-mailer/mailservice"
)

func main() {
	lambda.Start(mailservice.HandleLambdaEvent)
}
