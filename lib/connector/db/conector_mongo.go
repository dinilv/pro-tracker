package db

import "go.mongodb.org/mongo-driver/mongo"

//ENV variables
//load mongo essentials
type mongoDB struct {
	client mongo.Client
}

func newMongoDB() mongoDB {
	//TODO: mongo crdentials to be added to env
	return mongoDB{}
}

func (db mongoDB) Save(collection string, obj interface{}) error {
	//TODO: implement mongo specific ops
	return nil
}

func (db mongoDB) List(int, int, string, []string, map[string]interface{}, map[string]interface{}) ([]map[string]interface{}, error) {
	//TODO: implement mongo specific ops
	return []map[string]interface{}{}, nil
}
