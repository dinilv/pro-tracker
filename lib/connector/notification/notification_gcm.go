package notification

import (
	"context"
	"log"

	"google.golang.org/api/fcm/v1"
)

// google credentials
type googleCloudMessage struct {
	session fcm.Service
}

func newGoogleCloudMessage() googleCloudMessage {
	ctx := context.Background()
	fcmsess, err := fcm.NewService(ctx)
	if err != nil {
		log.Panic("GCM couldnt be connected.Existing...", err)
	}
	return googleCloudMessage{session: *fcmsess}
}

func (gcm googleCloudMessage) Send(string, string) error {
	//TODO: implement google messaging
	return nil
}
