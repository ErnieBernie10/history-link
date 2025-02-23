package record

import (
	data "historylink/internal/db"

	"github.com/google/uuid"
)

type recordResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	Significance string    `json:"significance"`
	Url          string    `json:"url"`
	StartDate    string    `json:"startDate"`
	EndDate      string    `json:"endDate"`
	Status       Status    `json:"status"`
	Type         Type      `json:"type"`
}

type createRecordCommandBody struct {
	Title        string `json:"title" validate:"required,max=255"`
	Description  string `json:"description" validate:"required,max=255"`
	Location     string `json:"location" validate:"max=255"`
	Significance string `json:"significance" validate:"max=255"`
	Url          string `json:"url" validate:"required,max=255"`
	StartDate    string `json:"startDate" validate:"required,datetime=200601021504"`
	EndDate      string `json:"endDate" validate:"required,datetime=200601021504"`
	Status       Status `json:"status" validate:"required,oneof=0 1 2 3"`
	Type         Type   `json:"type" validate:"required,oneof=0 1 2 3"`
}

func toResponse(record data.Record) recordResponse {
	return recordResponse{
		ID:           record.ID,
		Title:        record.Title,
		Description:  record.Description,
		Location:     record.Location.String,
		Significance: record.Significance.String,
		Url:          record.Url,
		StartDate:    record.StartDate.Time.Format("200601021504"),
		EndDate:      record.EndDate.Time.Format("200601021504"),
		Status:       Status(record.Status),
		Type:         Type(record.Type),
	}
}
