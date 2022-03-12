package mailgun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mlgun "github.com/mailgun/mailgun-go/v3"
	"gitlab.com/pro-tracker/lib/connector/bucket"
	"gitlab.com/pro-tracker/lib/connector/db"
	"gitlab.com/pro-tracker/lib/connector/notification"
	"gitlab.com/pro-tracker/lib/connector/queue"
	"gitlab.com/pro-tracker/tracker/constants"
	"gitlab.com/pro-tracker/tracker/model"
)

const (
	STORAGE_PARAMETER  = "STORAGE_PARAMETER"
	FOLDER_DATE_FORMAT = "02-01-2006"
	VERIFIED           = "verified"
	FAILED             = "failed"
	ROOT_FOLDER        = "pro-tracker"
	NEWLINE            = "\n"
)

// mailgun events are mapped to standard queue based on event type
var mailGunEventMapping map[string]string
var CALL_BACK_LOG string

func init() {
	// initialize queue based on env
	stage := os.Getenv("stage") // "dev" only
	clickQueue := "click-callback-" + stage + ".fifo"
	openQueue := "open-callback-" + stage + ".fifo"
	mailGunEventMapping = map[string]string{
		"clicked": clickQueue,
		"opened":  openQueue,
	}
	CALL_BACK_LOG = constants.CALL_BACK_LOG + stage
}

// 'Transaction Logic pattern'
// utility to carry out different services to callback handler
type CallbackServiceIF interface {
	VerifySignature(string, string, string) bool
	UploadMessage(string, []byte) error
	SendMessage(*model.Callback) error
	SendToQueue(*model.Callback) error
	SaveToStore(*model.Callback) error
}
type callbackServiceImpl struct {
	//private class variables for restricted access
	mailgunClient *mlgun.MailgunImpl
	bucketClient  bucket.BucketClient
	dbClient      db.DatabaseClient
	notifyClient  notification.NotificationClient
	queueClient   queue.QueueClient
}

// all dependencies are loaded based on the environment vars
func NewCallbackService() CallbackServiceIF {
	// take credentials from environment
	domain := os.Getenv(constants.DOMAIN)
	signKey := os.Getenv(constants.SIGN_KEY)
	mgClient := mlgun.NewMailgun(domain, signKey)
	// resolve bucket storage
	var bucketClient bucket.BucketClient
	fileStorageProvider := os.Getenv(constants.FILE_STORE_PROVIDER)
	bucketClient = bucket.New(fileStorageProvider)
	// resolve notification client
	var notificationClient notification.NotificationClient
	notificationProvider := os.Getenv(constants.NOTIFICATION_PROVIDER)
	notificationClient = notification.New(notificationProvider)
	// resolve database client
	var dbClient db.DatabaseClient
	dbProvider := os.Getenv(constants.STORAGE_PROVIDER)
	dbClient = db.New(dbProvider)
	// resolve queue client
	var queueClient queue.QueueClient
	qProvider := os.Getenv(constants.QUEUE_PROVIDER)
	queueClient = queue.New(qProvider)

	return callbackServiceImpl{
		mailgunClient: mgClient,
		bucketClient:  bucketClient,
		dbClient:      dbClient,
		notifyClient:  notificationClient,
		queueClient:   queueClient,
	}
}

// verify signature for event received
func (service callbackServiceImpl) VerifySignature(timestamp, token, signature string) bool {
	//use mailgun sdk to verify signature
	isValid, err := service.mailgunClient.VerifyWebhookSignature(mlgun.Signature{
		TimeStamp: timestamp,
		Token:     token,
		Signature: signature,
	})
	if err != nil {
		log.Println("Error in verifying mailgun event signature:", timestamp, err.Error())
		return false
	}
	return isValid
}

// help in uploading raw message to file storage
func (service callbackServiceImpl) UploadMessage(fileName string, fileBytes []byte) error {
	dir := "mailgun/"
	// generate folder name by date received
	dateFolder := time.Now().Format(FOLDER_DATE_FORMAT)
	dir = fmt.Sprintf("%s/%s", dir, dateFolder)
	err := service.bucketClient.Upload(ROOT_FOLDER, dir, fileName, fileBytes)
	if err != nil {
		log.Println("File upload to S3 failed.", err.Error())
	}
	return err
}

// send callback details for message.
func (service callbackServiceImpl) SendMessage(callback *model.Callback) error {
	// create message
	var buffer bytes.Buffer
	buffer.WriteString("Email: " + callback.Recipient + NEWLINE)
	buffer.WriteString("Provider: " + callback.Recipient + NEWLINE)
	buffer.WriteString("Action: " + callback.Event + NEWLINE)
	// format received at
	receivedAt := callback.ReceivedAt.Format("2006-01-02 15:04:05")
	buffer.WriteString("Received Time: " + receivedAt + NEWLINE)
	message := buffer.String()
	// send msg
	return service.notifyClient.Send("+919019138633", message)
}

// send a copy for co-relation and further processing
func (service callbackServiceImpl) SendToQueue(callback *model.Callback) error {
	// serialize log for sending to queue
	jsonBytes, err := json.Marshal(callback)
	if err != nil {
		log.Println("Couldnt serialize call back model.", err)
		return err
	}
	message := string(jsonBytes)
	// resolve queue for the event
	queueURL := mailGunEventMapping[callback.Event]
	// send to queue
	err = service.queueClient.Send(callback.MessageID, message, queueURL)
	return err
}

// keep a log of event in storage
func (service callbackServiceImpl) SaveToStore(callback *model.Callback) error {
	err := service.dbClient.Save(CALL_BACK_LOG, *callback)
	return err
}
