package queue

import (
	"log"

	"gitlab.com/pro-tracker/lib/constant"
)

// generic format for queue message
type Message struct {
	Body     *string
	HandleID *string
}

// signature for qClients
type QueueClient interface {
	Send(string, string, string) error
	ReceiveBatch(int, string) ([]Message, error)
	Delete(string, string) error
}

//  Factory pattern: according to vendor provide New QConnection
func New(category string) QueueClient {
	switch category {
	case constant.AMAZON_SQS:
		return newAmazonSQS()
	case constant.KAFKA:
		return newKafka()
	case constant.MOCK_QUEUE:
		return newMockQueue()
	default:
		log.Println("Couldnt find the required implementation for queue category:", category)
		return nil

	}
}
