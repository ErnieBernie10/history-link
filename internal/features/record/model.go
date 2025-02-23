package record

import (
	"github.com/google/uuid"
)

type pagedResponse[t any] struct {
	Page    int `json:"page"`
	Size    int `json:"size"`
	Total   int `json:"total"`
	Records []t `json:"records"`
}

type impactResponse struct {
	Value       int       `json:"value"`
	Category    int       `json:"category"`
	Description string    `json:"description"`
	ID          uuid.UUID `json:"id"`
	RecordID    uuid.UUID `json:"recordId"`
}

type recordResponseBody struct {
	ID           uuid.UUID        `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Location     string           `json:"location"`
	Significance string           `json:"significance"`
	Url          string           `json:"url"`
	StartDate    string           `json:"startDate"`
	EndDate      string           `json:"endDate"`
	RecordStatus RecordStatus     `json:"recordStatus"`
	Type         Type             `json:"type"`
	Impacts      []impactResponse `json:"impacts"`
}

type createRecordCommandBody struct {
	Title        string                    `json:"title" required:"true" maxLength:"255"`
	Description  string                    `json:"description" validate:"required,max=255"`
	Location     string                    `json:"location" validate:"max=255"`
	Significance string                    `json:"significance" validate:"max=255"`
	Url          string                    `json:"url" validate:"required,max=255"`
	StartDate    string                    `json:"startDate" validate:"required,datetime=200601021504"`
	EndDate      string                    `json:"endDate" validate:"required,datetime=200601021504"`
	RecordStatus RecordStatus              `json:"recordStatus" validate:"required,oneof=0 1 2 3"`
	Type         Type                      `json:"type" validate:"required,oneof=0 1 2 3"`
	Impacts      []createImpactCommandBody `json:"impacts" required:"true"`
}

type updateRecordCommandBody struct {
	ID uuid.UUID `json:"id" path:"id"`

	Title        string                    `json:"title" validate:"max=255"`
	Description  string                    `json:"description" validate:"max=255"`
	Location     string                    `json:"location" validate:"max=255"`
	Significance string                    `json:"significance" validate:"max=255"`
	Url          string                    `json:"url" validate:"required,max=255"`
	StartDate    string                    `json:"startDate" validate:"required,datetime=200601021504"`
	EndDate      string                    `json:"endDate" validate:"required,datetime=200601021504"`
	RecordStatus RecordStatus              `json:"recordStatus" validate:"required,oneof=0 1 2 3"`
	Type         Type                      `json:"type" validate:"required,oneof=0 1 2 3"`
	Impacts      []updateImpactCommandBody `json:"impacts" required:"true"`
}

type createImpactCommandBody struct {
	Description string `json:"description"`
	Value       int    `json:"value"`
	Category    int    `json:"category"`
}

type updateImpactCommandBody struct {
	ID          uuid.UUID `json:"id" path:"id"`
	Description string    `json:"description"`
	Value       int       `json:"value"`
	Category    int       `json:"category"`
	RecordId    uuid.UUID `json:"recordId"`
}

func toResponse(record Record) recordResponseBody {
	return recordResponseBody{
		ID:           record.ID,
		Title:        record.Title,
		Description:  record.Description,
		Location:     record.Location.String,
		Significance: record.Significance.String,
		Url:          record.Url,
		StartDate:    record.StartDate.Time.Format("200601021504"),
		EndDate:      record.EndDate.Time.Format("200601021504"),
		RecordStatus: RecordStatus(record.RecordStatus),
		Type:         Type(record.Type),
		Impacts:      toImpactResponse(record.Impacts),
	}
}

func toImpactResponse(impacts []Impact) []impactResponse {
	var response []impactResponse
	for _, impact := range impacts {
		response = append(response, impactResponse{
			Value:       impact.Value,
			Category:    impact.Category,
			Description: impact.Description,
			ID:          impact.ID,
			RecordID:    impact.RecordID,
		})
	}
	return response
}

func toPagedResponse(records []Record, page int, pageSize int) pagedResponse[recordResponseBody] {
	var recordsResponse []recordResponseBody
	for _, record := range records {
		recordsResponse = append(recordsResponse, toResponse(record))
	}
	return pagedResponse[recordResponseBody]{
		Page:    page,
		Size:    pageSize,
		Total:   len(records),
		Records: recordsResponse,
	}
}
