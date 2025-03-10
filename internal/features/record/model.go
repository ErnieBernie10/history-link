package record

import (
	"historylink/internal/common"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type pagedResponse[t any] struct {
	Page    int `json:"page"`
	Size    int `json:"size"`
	Total   int `json:"total"`
	Records []t `json:"records"`
}

type impactResponse struct {
	Value       int16     `json:"value"`
	Category    Category  `json:"category"`
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
	UpdatedAt    string           `json:"updatedAt"`
	CreatedAt    string           `json:"createdAt"`
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
	RecordStatus RecordStatus              `json:"recordStatus" enum:"removed,draft,pending,reviewed"`
	Type         Type                      `json:"type" enum:"arc,event,person,object"`
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
	RecordStatus RecordStatus              `json:"recordStatus" enum:"removed,draft,pending,reviewed"`
	Type         Type                      `json:"type" enum:"arc,event,person,object"`
	Impacts      []updateImpactCommandBody `json:"impacts"`
}

type createImpactCommandBody struct {
	Description string   `json:"description" minLength:"1" maxLength:"255"`
	Value       int16    `json:"value" minimum:"1" maximum:"10"`
	Category    Category `json:"category" enum:"economic,political,social,cultural,tech"`
}

type updateImpactCommandBody struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Description string    `json:"description" minLength:"1" maxLength:"255"`
	Value       int16     `json:"value" minimum:"1" maximum:"10"`
	Category    Category  `json:"category" enum:"economic,political,social,cultural,tech"`
	RecordId    uuid.UUID `json:"recordId,omitempty"`
}

func (record RecordAggregate) toResponse() recordResponseBody {
	return recordResponseBody{
		ID:           record.ID,
		Title:        record.Title,
		Description:  record.Description,
		Location:     record.Location,
		Significance: record.Significance,
		Url:          record.URL,
		StartDate:    common.ToDateString(record.StartDate),
		EndDate:      common.ToDateString(record.EndDate),
		RecordStatus: RecordStatusFromInt16(record.Status),
		Type:         TypeFromInt16(record.Type),
		Impacts: lo.Map(record.Impacts, func(impact ImpactEntity, index int) impactResponse {
			return impact.toResponse()
		}),
		UpdatedAt: common.ToDateTimeString(&record.History.UpdatedAt),
		CreatedAt: common.ToDateTimeString(&record.History.CreatedAt),
	}
}

func (i ImpactEntity) toResponse() impactResponse {
	return impactResponse{
		ID:          i.ID,
		Value:       i.Value,
		Category:    CategoryFromInt16(i.Category),
		Description: i.Description,
		RecordID:    i.RecordID,
	}
}
