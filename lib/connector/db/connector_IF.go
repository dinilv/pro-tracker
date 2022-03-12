package db

import (
	"log"

	"gitlab.com/pro-tracker/lib/constant"
)

// factory for storing data & logs
type DatabaseClient interface {
	Save(string, interface{}) error
	List(int, int, string, []string, map[string]interface{}, map[string]interface{}) ([]map[string]interface{}, error)
}

func New(category string) DatabaseClient {
	switch category {
	case constant.DYNAMO_DB:
		return newDynamoDB()
	case constant.MONGO_DB:
		return newMongoDB()
	case constant.MOCK_DB:
		return newMockDB()
	default:
		log.Println("Couldnt find the required implementation for db category:", category)
		return nil
	}
}
