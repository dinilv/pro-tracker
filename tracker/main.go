package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.com/pro-tracker/tracker/handler/mailgun"
)

func main() {
	log.Printf("Tracker Server is starting UP")
	//inititate services for handlers
	mailgunService := mailgun.NewCallbackService() // will initiate dependencies according to ENV
	mailgunHandler := mailgun.NewCallbackHandler(mailgunService)
	lambda.Start(mailgunHandler.Callback)
}
