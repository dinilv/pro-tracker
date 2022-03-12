package notification

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Provide new AmazonSNS handler implementation
// Send(): message to a phone
type amazonSNS struct {
	session *sns.SNS
}

// ENV vars
// AWS_ACCESS_KEY_ID=XXXXXXXXXX
// AWS_SECRET_ACCESS_KEY=XXXXXXXX
// AWS_REGION=XXX
func newAmazonSNS() amazonSNS {
	sess := session.Must(session.NewSession())
	snsSess := sns.New(sess)
	return amazonSNS{session: snsSess}
}

func (amazonSns amazonSNS) Send(phone, message string) error {
	//create params
	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phone),
	}
	resp, err := amazonSns.session.Publish(params)
	//handle error if any
	if err != nil {
		log.Println("Message sending failed.", err)
		return err
	}
	//pretty print response
	log.Println("Message Sent successfully.", resp)
	return nil
}
