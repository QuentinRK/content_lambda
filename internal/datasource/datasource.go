package datasource

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DB interface {
	*dynamodb.Client | string
}

type DataSource interface {
	FetchByRecordType(ctx context.Context, recordType string) ([]WebsiteRecord, error)
}

type WebsiteRecord struct {
	RecordType  string `dynamodbav:recordType`
	Name        string `dynamodbav:name`
	Date        string `dynamodbav:date`
	Description string `dynamodbav:description`
	Link        string `dynamodbav:link`
	Location    string `dynamodbav:location`
	OrderNo     int    `dynamodbav:orderNo`
	Role        string `dynamodbav:role`
}
