package notification

import (
	"log"

	"gitlab.com/pro-tracker/lib/constant"
)

// Factory for sending notification
type NotificationClient interface {
	Send(string, string) error
}

func New(category string) NotificationClient {
	switch category {
	case constant.AMAZON_SNS:
		return newAmazonSNS()
	case constant.GOOGLE_CLOUD_MESSAGE:
		return newGoogleCloudMessage()
	case constant.MOCK_NOTIFICATION:
		return newMockNotification()
	default:
		log.Println("Couldnt find the required implementation for notification category:", category)
		return nil
	}
}
