package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Provide new DynamoDB handler implementation
// Save(): for saving a record
// List(): query for multiple records
type dynamoDB struct {
	session *dynamodb.DynamoDB
}

// ENV vars
// AWS_ACCESS_KEY_ID=XXXXXXXXXX
// AWS_SECRET_ACCESS_KEY=XXXXXXXX
// AWS_REGION=XXX
func newDynamoDB() DatabaseClient {
	sess := session.Must(session.NewSession())
	dynamoCli := dynamodb.New(sess)
	return dynamoDB{session: dynamoCli}
}

// save record to table
func (db dynamoDB) Save(table string, obj interface{}) error {
	//marshal to amazon value object
	av, err := dynamodbattribute.MarshalMap(obj)
	if err != nil {
		log.Println("Got error marshalling new save Obj:", err)
	}
	// save to dynamoDB
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}
	_, err = db.session.PutItem(input)
	if err != nil {
		log.Println("Got error calling PutItem:", err)
	}
	return err
}

// list items from table
func (db dynamoDB) List(page int, limit int, table string, fields []string, order map[string]interface{}, filter map[string]interface{}) ([]map[string]interface{}, error) {
	var nameFilter expression.ConditionBuilder
	//support multiple filter conditions
	for k, v := range filter {
		nameFilter = expression.Name(k).Equal(expression.Value(v))
	}
	//selecting required fields
	var proj expression.ProjectionBuilder
	for _, field := range fields {
		proj = expression.NamesList(expression.Name(field))
	}
	//build expression to query
	expr, err := expression.NewBuilder().WithFilter(expression.ConditionBuilder(nameFilter)).WithProjection(proj).Build()
	if err != nil {
		log.Println("Got error building expression:", err)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table),
	}
	result, err := db.session.Scan(params)
	if err != nil {
		log.Println("Query DynamoDB call failed:", err)
		return nil, err
	}
	// transform items to generic format
	items := []map[string]interface{}{}
	for _, i := range result.Items {
		item := map[string]interface{}{}
		items = append(items, item)
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Println("Got error unmarshalling:", err)
			break
		}
	}
	return items, err

}
