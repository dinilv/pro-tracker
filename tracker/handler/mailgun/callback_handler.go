package mailgun

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/pro-tracker/tracker/model"
)

const (
	TIMESTAMP = "timestamp"
	BODY      = "body"
	SIGNATURE = "signature"
	TOKEN     = "token"
)

type CallbackHandlerIF interface {
	Callback(context.Context, json.RawMessage) (events.APIGatewayProxyResponse, error)
}

// 'Manager Pattern' on handler
// parse the event and handle api responses
// will only delegate transactions to services
type CallbackHandler struct {
	service CallbackServiceIF
}

func NewCallbackHandler(service CallbackServiceIF) CallbackHandlerIF {
	return CallbackHandler{service: service}
}

func (handler CallbackHandler) Callback(ctx context.Context, data json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// declare an empty map interface to unmarshall req
	var event map[string]interface{}
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(data, &event)

	// extract signature details of event & message body
	timestamp, token, signature, eventBody, err := parse(event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401, // unauthenticated error
		}, err
	}
	// validate signature of messge by custom validation
	isVerified := handler.service.VerifySignature(timestamp, token, signature)
	if !isVerified {
		return events.APIGatewayProxyResponse{
			StatusCode: 401, // unauthenticated error
		}, errors.New("signature verification failed")
	}

	// parse call back event to object model
	callbackObj := &model.Callback{}
	callbackObj.ParseMailgun(eventBody)

	// sending SMS regarding the client is priority Business Operation.
	err = handler.service.SendMessage(callbackObj)
	if err != nil {
		// respond with error-code 406 to have retry from mailgun later if SMS failed
		return events.APIGatewayProxyResponse{
			StatusCode: 406,
		}, errors.New("sending notification failed")
	}
	callbackObj.IsMessaged = true

	// save to appropriate S3 bucket
	err = handler.service.UploadMessage(callbackObj.MessageID, data)
	if err == nil {
		// upload to s3 suceeded, but even if its failed will continue in path to send SMS
		callbackObj.IsFiled = true
	}

	// push to queue for further processing,co-relation & analysis
	err = handler.service.SendToQueue(callbackObj)
	if err == nil {
		// err is not a blocking business operation for proceeding
		callbackObj.IsQueued = true
	}

	// save the log to DynamoDB
	err = handler.service.SaveToStore(callbackObj)
	if err != nil {
		// err is not a blocking business operation for proceeding to sucess response
		log.Println("Save to DB store failed", err)
	}
	// respond with success
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Webhook received from mailgun has been sucessfully processed",
	}, nil
}

func parse(event map[string]interface{}) (timestamp, token, signature string, body map[string]interface{}, err error) {
	bodyInterface, isBody := event[BODY]
	if !isBody || bodyInterface == nil {
		// body not available part of invocation: failed to read body
		err = errors.New("body for event not available")
		return
	}
	bodyText := bodyInterface.(string)
	body = map[string]interface{}{}
	err = json.Unmarshal([]byte(bodyText), &body)
	if err != nil {
		err = errors.New("not able to unmarshall body mesage")
		return
	}
	// get text body & convert to json
	signatureInteface, isSignatureBody := body[SIGNATURE]
	if !isSignatureBody {
		// signature not available part of event: failed validation for mailgun
		err = errors.New("signature for event not available")
		return
	}
	signatureMap := signatureInteface.(map[string]interface{}) // type case to concrete type
	timestampIF, isTimestamp := signatureMap[TIMESTAMP]
	tokenIF, isToken := signatureMap[TOKEN]
	signatureIF, isSignature := signatureMap[SIGNATURE]
	if !isTimestamp || !isSignature || !isToken {
		// necessary entry missing
		err = errors.New("signature for event not available")
		return
	}
	timestamp = timestampIF.(string)
	token = tokenIF.(string)
	signature = signatureIF.(string)

	for k, v := range body {
		log.Println("Body Key Value", k, v)
	}
	return
}
