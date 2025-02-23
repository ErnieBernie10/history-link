package record

import (
	"context"
	"errors"
	"historylink/internal/db"
	"time"

	"github.com/google/uuid"
)

func NewRecordService(recordRepository IRecordRepository) IRecordService {
	return RecordService{
		recordRepository: recordRepository,
	}
}

type IRecordService interface {
	Create(c context.Context, command createRecordCommandBody) (Record, error)
	Update(c context.Context, id uuid.UUID, command updateRecordCommandBody) error
	GetById(id uuid.UUID) (Record, error)
	GetPaged(c context.Context, page, pageSize int) ([]Record, error)
}

type RecordService struct {
	recordRepository IRecordRepository
}

type RecordStatus int16

const (
	Removed RecordStatus = iota
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

func (s RecordService) Create(context context.Context, command createRecordCommandBody) (Record, error) {
	startDate, err := time.Parse("200601021504", command.StartDate)
	if err != nil {
		return Record{}, err
	}
	endDate, err := time.Parse("200601021504", command.EndDate)
	if err != nil {
		return Record{}, err
	}
	return s.recordRepository.Create(context, Record{
		Title:        command.Title,
		Description:  command.Description,
		Location:     db.NewNullString(command.Location),
		Significance: db.NewNullString(command.Significance),
		Url:          command.Url,
		StartDate:    db.NewNullTime(startDate),
		EndDate:      db.NewNullTime(endDate),
		Type:         int16(command.Type),
		RecordStatus: int16(command.RecordStatus),
		Impacts:      mapCreateImpacts(command.Impacts),
	})
}

func mapUpdateImpacts(impactCommands []updateImpactCommandBody) []Impact {
	impacts := make([]Impact, len(impactCommands))
	for i, command := range impactCommands {
		impacts[i] = Impact{
			ID:          command.ID,
			Description: command.Description,
			Value:       command.Value,
			Category:    command.Category,
			RecordID:    command.RecordId,
		}
	}
	return impacts
}

func mapCreateImpacts(impactCommands []createImpactCommandBody) []Impact {
	impacts := make([]Impact, len(impactCommands))
	for i, command := range impactCommands {
		impacts[i] = Impact{
			Description: command.Description,
			Value:       command.Value,
			Category:    command.Category,
		}
	}
	return impacts
}

func (s RecordService) GetById(id uuid.UUID) (Record, error) {
	return s.recordRepository.GetById(id)
}

func (s RecordService) Update(c context.Context, id uuid.UUID, command updateRecordCommandBody) error {
	if id != command.ID {
		return errors.New("id mismatch")
	}
	startDate, err := time.Parse("200601021504", command.StartDate)
	if err != nil {
		return err
	}
	endDate, err := time.Parse("200601021504", command.EndDate)
	if err != nil {
		return err
	}
	return s.recordRepository.Update(c, Record{
		ID:           command.ID,
		Title:        command.Title,
		Description:  command.Description,
		Location:     db.NewNullString(command.Location),
		Significance: db.NewNullString(command.Significance),
		Url:          command.Url,
		StartDate:    db.NewNullTime(startDate),
		EndDate:      db.NewNullTime(endDate),
		Type:         int16(command.Type),
		RecordStatus: int16(command.RecordStatus),
		Impacts:      mapUpdateImpacts(command.Impacts),
	})
}

func (s RecordService) GetPaged(c context.Context, page int, pageSize int) ([]Record, error) {
	return s.recordRepository.GetPaged(c, pageSize, (page-1)*pageSize)
}
