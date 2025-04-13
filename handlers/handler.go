package handlers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/quentinrankin/content_lambda/internal/datasource"
	"github.com/quentinrankin/content_lambda/internal/repository"
)

const (
	defaultTableName = "websiteData"
	pathAboutMe      = "/about-me"
	pathProjects     = "/projects"
	pathWorkHistory  = "/work"
)

func Init() func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(GetEnv("AWS_REGION", "eu-west-1")))
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(cfg)

	tableName := GetEnv("WEBSITE_TABLE_NAME", "portfolio_content")

	repo := repository.WebsiteRepository{
		DS: &datasource.DynamodbDataSource{
			TableName: tableName,
			DB:        dbClient,
		},
	}

	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		var content string
		var err error
		contentPath := request.RequestContext.HTTP.Path

		switch contentPath {
		case pathAboutMe:
			records, fetchErr := repo.GetAboutMe(ctx)
			if fetchErr != nil {
				return handleError(fetchErr, "fetching about-me data")
			}
			content, err = repo.GetAboutMeResponse(records)

		case pathProjects:
			records, fetchErr := repo.GetProjects(ctx)
			if fetchErr != nil {
				return handleError(fetchErr, "fetching projects data")
			}
			content, err = repo.GetProjectsResponse(records)

		case pathWorkHistory:
			records, fetchErr := repo.GetWorkHistory(ctx)
			if fetchErr != nil {
				return handleError(fetchErr, "fetching work history data")
			}
			content, err = repo.GetWorkHistoryResponse(records)

		default:
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       fmt.Sprintf(`{"error": "Route %s not found"}`, contentPath),
			}, nil
		}

		if err != nil {
			return handleError(err, "processing response")
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       content,
		}, nil
	}
}

func handleError(err error, context string) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("Error %s: %v", context, err)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf(`{"error": "An error occurred while %s"}`, context),
	}, nil
}

func GetEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
