package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/quentinrankin/content_lambda/internal/datasource"
)

type WebsiteRepository struct {
	DS datasource.DataSource
}

type RecordType string

const (
	RecordTypeBio      RecordType = "biography"
	RecordTypeProjects RecordType = "project"
	RecordTypeWork     RecordType = "work"
)

type Response interface {
	AboutMeResponse | ProjectResponse | []ProjectResponse | []WorkHistoryResponse
}

type AboutMeResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type WorkHistoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Location    string `json:"location"`
	Date        string `json:"date"`
	OrderNo     int    `json:"orderNo"`
}

func (wr *WebsiteRepository) fetchRecords(ctx context.Context, recordType RecordType) ([]datasource.WebsiteRecord, error) {
	results, err := wr.DS.FetchByRecordType(ctx, string(recordType))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s records: %w", recordType, err)
	}

	return results, nil
}

func (wr *WebsiteRepository) GetAboutMe(ctx context.Context) ([]datasource.WebsiteRecord, error) {
	return wr.fetchRecords(ctx, RecordTypeBio)
}

func (wr *WebsiteRepository) GetProjects(ctx context.Context) ([]datasource.WebsiteRecord, error) {
	return wr.fetchRecords(ctx, RecordTypeProjects)
}

func (wr *WebsiteRepository) GetWorkHistory(ctx context.Context) ([]datasource.WebsiteRecord, error) {
	return wr.fetchRecords(ctx, RecordTypeWork)
}

func (wr *WebsiteRepository) GetAboutMeResponse(records []datasource.WebsiteRecord) (string, error) {
	if len(records) == 0 {
		return "[]", nil
	}

	record := records[0]
	return toJSON(AboutMeResponse{
		Name:        record.Name,
		Description: record.Description,
	})
}

func (wr *WebsiteRepository) GetProjectsResponse(records []datasource.WebsiteRecord) (string, error) {
	response := []ProjectResponse{}

	for _, record := range records {
		response = append(response, ProjectResponse{
			Name:        record.Name,
			Description: record.Description,
			Link:        record.Link,
		})
	}

	return toJSON(response)
}

func (wr *WebsiteRepository) GetWorkHistoryResponse(records []datasource.WebsiteRecord) (string, error) {
	response := []WorkHistoryResponse{}

	for _, record := range records {
		response = append(response, WorkHistoryResponse{
			Name:        record.Name,
			Description: record.Description,
			Location:    record.Location,
			Date:        record.Date,
			OrderNo:     record.OrderNo,
			Role:        record.Role,
		})
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].OrderNo < response[j].OrderNo
	})

	return toJSON(response)
}

func toJSON[T Response](data T) (string, error) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	return string(jsonResponse), nil
}
