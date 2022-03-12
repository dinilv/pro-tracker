package bucket

import (
	"log"

	"gitlab.com/pro-tracker/lib/constant"
)

//Bucket Factory for file storage.
type BucketClient interface {
	Upload(string, string, string, []byte) error
}

func New(category string) BucketClient {
	switch category {
	case constant.AMAZON_STORAGE:
		return newAmazonStorage()
	case constant.FILESTACK:
		return newFilestack()
	case constant.MOCK_FILE_STORE:
		return newBucketMock()
	default:
		log.Println("Couldnt find the required implementation for bucket category:", category)
		return nil
	}
}
