package datasource

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamodbDataSource struct {
	TableName string
	DB        *dynamodb.Client
}

func GetExpression(exprMap map[string]string) expression.Expression {
	var expr expression.Expression
	for key, value := range exprMap {
		keyExpr := expression.Key(key).Equal(expression.Value(value))

		var err error
		expr, err = expression.NewBuilder().WithKeyCondition(keyExpr).Build()
		if err != nil {
			panic(err)
		}

	}

	return expr
}

func (DS *DynamodbDataSource) FetchByRecordType(ctx context.Context, recordType string) ([]WebsiteRecord, error) {
	var results []WebsiteRecord

	expr := GetExpression(map[string]string{
		"recordType": recordType,
	})

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(DS.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	queryPaginator := dynamodb.NewQueryPaginator(DS.DB, input)
	for queryPaginator.HasMorePages() {
		page, err := queryPaginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}

		var pageItems []WebsiteRecord
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &pageItems); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		results = append(results, pageItems...)
	}

	return results, nil
}
