package bucket

import (
	"bytes"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	BUCKET_ACCESS  = "private"
	BUCKET_ENCRYPT = "AES256"
	ATTACHMENT     = "attachment"
)

// Provide new S3 handler implementation
// Upload(): for storing a file
// private access variables
type amazonStorage struct {
	client *s3.S3
}

// Env variables for amz
// AWS_ACCESS_KEY_ID=XXXXXXXXXX
// AWS_SECRET_ACCESS_KEY=XXXXXXXX
// AWS_DEFAULT_REGION=xxx
func newAmazonStorage() amazonStorage {
	session, err := session.NewSession()
	if err != nil {
		log.Fatalln("Failed to initiate amazon session.", err, session)
	}
	s3 := s3.New(session)
	return amazonStorage{client: s3}
}

// store file in bucket
func (store amazonStorage) Upload(bucket, folder, fileName string, fileBytes []byte) error {
	_, err := store.client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(folder + "/" + fileName + ".txt"),
		ACL:                  aws.String(BUCKET_ACCESS),
		Body:                 bytes.NewReader(fileBytes),
		ContentLength:        aws.Int64(int64(len(fileBytes))),
		ContentType:          aws.String(http.DetectContentType(fileBytes)),
		ContentDisposition:   aws.String(ATTACHMENT),
		ServerSideEncryption: aws.String(BUCKET_ENCRYPT),
	})
	if err != nil {
		log.Println("Coudn't upload to s3 bucket.", bucket, folder, err)
		return err
	}
	log.Println("Uploaded file to s3 bucket.", bucket, folder)
	return nil
}
