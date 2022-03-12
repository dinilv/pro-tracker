package constant

const (
	//database
	DYNAMO_DB = "dynamo_db"
	MONGO_DB  = "mongo_db"
	MOCK_DB   = "mock_db"
	//file bucket
	AMAZON_STORAGE  = "amazon_s3"
	FILESTACK       = "filestack"
	MOCK_FILE_STORE = "mock_file_store"
	//env vars
	AWS_S3_REGION         = "AMAZON_S3_REGION"
	AWS_S3_MAILGUN_BUCKET = "pro-tracker/webhook/mailgun"
	//queue
	AMAZON_SQS = "amazon_sqs"
	KAFKA      = "kafka"
	MOCK_QUEUE = "mock_queue"
	//notification
	AMAZON_SNS           = "amazon_sns"
	GOOGLE_CLOUD_MESSAGE = "google_cloud_message"
	MOCK_NOTIFICATION    = "mock_notification"
)
