package queue

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Env variables to be set
// AWS_SDK_LOAD_CONFIG=1
// AWS_ACCESS_KEY_ID=XXXXXXXXXX
// AWS_SECRET_ACCESS_KEY=XXXXXXXX
type amazonSQS struct {
	client *sqs.SQS
}

// Provide new SQS handler implementation
// Send(): for adding to queue
// Delete(): after processing the msg
// ReceiveBatch(): pull message in batch
func newAmazonSQS() amazonSQS {
	// create amz session using env creds
	session, err := session.NewSession()
	if err != nil {
		log.Panic("Failed to initiate AmazonSQS", err, session)
	}
	//create new queue client
	sqsClient := sqs.New(session)
	return amazonSQS{client: sqsClient}
}

// send a msg
func (queue amazonSQS) Send(messageID, message, url string) error {
	// transform to AWS input ops
	out, err := queue.client.SendMessage(&sqs.SendMessageInput{
		QueueUrl:               &url,
		MessageBody:            aws.String(message),
		MessageGroupId:         &messageID,
		MessageDeduplicationId: &messageID,
	})
	if err != nil {
		log.Println("Sending to queue failed.", url, err)
	} else {
		log.Println("Sent message to queue sucessfullly.", url, out)
	}
	return err
}

// delete msg after processing its contents
func (queue amazonSQS) Delete(handleID, url string) error {
	out, err := queue.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: aws.String(handleID),
	})
	if err != nil {
		log.Println("Deleting message from queue failed.", url, err)
	} else {
		log.Println("Deleted from queue sucessfullly.", url, out)
	}
	return err
}

// pull msgs based on the count
func (queue amazonSQS) ReceiveBatch(batch int, url string) ([]Message, error) {
	msgResult, err := queue.client.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &url,
		MaxNumberOfMessages: aws.Int64(int64(batch)),
	})
	if err != nil {
		log.Println("Receiving message from queue failed.", url, err)
		return nil, err
	}
	// transform received message to standard format
	messages := []Message{}
	for _, msg := range msgResult.Messages {
		eachMsg := Message{
			Body:     msg.Body,
			HandleID: msg.ReceiptHandle,
		}
		messages = append(messages, eachMsg)
	}
	return messages, nil
}
