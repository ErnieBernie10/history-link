package record

import (
	"context"
	"database/sql"
	data "historylink/internal/db"
	"time"

	"github.com/google/uuid"
)

func NewRecordService(db *data.Queries) IRecordService {
	return RecordService{
		db: db,
	}
}

type IRecordService interface {
	Create(command createRecordCommand) (data.Record, error)
	GetById(id uuid.UUID) (data.Record, error)
}

type RecordService struct {
	db *data.Queries
}

type Status int16

const (
	Removed Status = iota
	Draft
	PendingReview
	Reviewed
)

type Type int16

const (
	Arc Type = iota
	Event
	Person
	Object
)

type createRecordCommand struct {
	Title        string
	Description  string
	Location     string
	Significance string
	Url          string
	StartDate    time.Time
	EndDate      time.Time
	Status       Status
	Type         Type
}

func (s RecordService) Create(command createRecordCommand) (data.Record, error) {
	return s.db.CreateRecord(context.Background(), data.CreateRecordParams{
		Title:        command.Title,
		Description:  command.Description,
		Location:     sql.NullString{String: command.Location, Valid: true},
		Significance: sql.NullString{String: command.Significance, Valid: true},
		StartDate:    sql.NullTime{Time: command.StartDate, Valid: true},
		EndDate:      sql.NullTime{Time: command.EndDate, Valid: true},
		Status:       int16(command.Status),
		Url:          command.Url,
		Type:         int16(command.Type),
	})
}

func (s RecordService) GetById(id uuid.UUID) (data.Record, error) {
	return s.db.GetRecord(context.Background(), id)
}
