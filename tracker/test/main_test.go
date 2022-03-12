package test

import (
	"context"
	"encoding/json"
	lib "gitlab.com/pro-tracker/lib/constant"
	"gitlab.com/pro-tracker/tracker/constants"
	"gitlab.com/pro-tracker/tracker/handler/mailgun"
	"os"
	"testing"
)

var mailgunService mailgun.CallbackServiceIF
var mailgunHandler mailgun.CallbackHandlerIF

func init() {
	// set env variables to use MOCK instances
	os.Setenv(constants.STORAGE_PROVIDER, lib.MOCK_DB)
	os.Setenv(constants.QUEUE_PROVIDER, lib.MOCK_QUEUE)
	os.Setenv(constants.NOTIFICATION_PROVIDER, lib.MOCK_NOTIFICATION)
	os.Setenv(constants.FILE_STORE_PROVIDER, lib.MOCK_FILE_STORE)
	// inititate services for handlers
	mailgunService = mailgun.NewCallbackService()
	mailgunHandler = mailgun.NewCallbackHandler(mailgunService)

}

func TestCallbackHandlerValidation(t *testing.T) {
	ctx := context.Background()
	// validation for error check
	_, err := mailgunHandler.Callback(ctx, json.RawMessage(""))
	// expecting error due to missing body text
	if err == nil {
		t.Errorf("Expected validation error not empty")
	}
}

// todo: create raw message and send for testing
func TestCallbackServiceParsing(t *testing.T) {
}
