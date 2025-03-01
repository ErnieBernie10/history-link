package record

import (
	"historylink/internal/common"

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
	Location     *string          `json:"location"`
	Significance *string          `json:"significance"`
	Url          string           `json:"url"`
	StartDate    string           `json:"startDate"`
	EndDate      string           `json:"endDate"`
	RecordStatus RecordStatus     `json:"recordStatus"`
	Type         Type             `json:"type"`
	Impacts      []impactResponse `json:"impacts"`
}

type createRecordCommandBody struct {
	Title        string                    `json:"title" minLength:"1" maxLength:"255"`
	Description  string                    `json:"description" minLength:"1" maxLength:"255"`
	Location     string                    `json:"location" minLength:"1" maxLength:"255"`
	Significance string                    `json:"significance" minLength:"1" maxLength:"255"`
	Url          string                    `json:"url" minLength:"1" maxLength:"255"`
	StartDate    string                    `json:"startDate" format:"date"`
	EndDate      string                    `json:"endDate" format:"date"`
	RecordStatus RecordStatus              `json:"recordStatus" enum:"0,1,2,3"`
	Type         Type                      `json:"type" enum:"0,1,2,3"`
	Impacts      []createImpactCommandBody `json:"impacts"`
}

type updateRecordCommandBody struct {
	ID           uuid.UUID                 `json:"id" path:"id"`
	Title        string                    `json:"title" minLength:"1" maxLength:"255"`
	Description  string                    `json:"description" minLength:"1" maxLength:"255"`
	Location     string                    `json:"location" minLength:"1" maxLength:"255"`
	Significance string                    `json:"significance" minLength:"1" maxLength:"255"`
	Url          string                    `json:"url" minLength:"1" maxLength:"255"`
	StartDate    string                    `json:"startDate" format:"date"`
	EndDate      string                    `json:"endDate" format:"date"`
	RecordStatus RecordStatus              `json:"recordStatus" enum:"0,1,2,3"`
	Type         Type                      `json:"type" enum:"0,1,2,3"`
	Impacts      []updateImpactCommandBody `json:"impacts"`
}

type createImpactCommandBody struct {
	Description string `json:"description" minLength:"1" maxLength:"255"`
	Value       int    `json:"value" minimum:"1" maximum:"10"`
	Category    int    `json:"category" minimum:"1"`
}

type updateImpactCommandBody struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Description string    `json:"description" minLength:"1" maxLength:"255"`
	Value       int       `json:"value" minimum:"1" maximum:"10"`
	Category    int       `json:"category" minimum:"1"`
	RecordId    uuid.UUID `json:"recordId,omitempty"`
}

func toResponse(record RecordAggregate) recordResponseBody {
	return recordResponseBody{
		ID:           record.ID,
		Title:        record.Title,
		Description:  record.Description,
		Location:     record.Location,
		Significance: record.Significance,
		Url:          record.URL,
		StartDate:    common.ToDateString(record.StartDate),
		EndDate:      common.ToDateString(record.EndDate),
		RecordStatus: RecordStatus(record.Status),
		Type:         Type(record.Type),
		Impacts:      toImpactResponse(record),
	}
}

func toImpactResponse(record RecordAggregate) []impactResponse {
	var response []impactResponse
	for _, impact := range record.Impacts {
		response = append(response, impactResponse{
			Value:       int(impact.Value),
			Category:    int(impact.Category),
			Description: impact.Description,
			ID:          impact.ID,
			RecordID:    *impact.RecordID,
		})
	}
	return response
}
