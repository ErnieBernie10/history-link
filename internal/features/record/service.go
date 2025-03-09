package record

import (
	"context"
	"errors"
	"historylink/.gen/historylink/public/model"
	"historylink/internal/common"
	"log/slog"

	"github.com/google/uuid"
)

func NewRecordService(recordRepository IRecordRepository, logger *slog.Logger) IRecordService {
	return RecordService{
		recordRepository: recordRepository,
		logger:           logger,
	}
}

type IRecordService interface {
	Create(c context.Context, command createRecordCommandBody) (recordResponseBody, error)
	Update(c context.Context, id uuid.UUID, command updateRecordCommandBody) error
	GetById(id uuid.UUID) (recordResponseBody, error)
	GetPaged(c context.Context, page, pageSize int) ([]recordResponseBody, int, error)
	Delete(c context.Context, id uuid.UUID) error
}

type RecordService struct {
	recordRepository IRecordRepository
	logger           *slog.Logger
}

type RecordStatus string

const (
	Removed       RecordStatus = "removed"
	Draft         RecordStatus = "draft"
	PendingReview RecordStatus = "pending"
	Reviewed      RecordStatus = "reviewed"
)

type Type string

const (
	Arc    Type = "arc"
	Event  Type = "event"
	Person Type = "person"
	Object Type = "object"
)

type Category string

const (
	Political Category = "political"
	Social    Category = "social"
	Economic  Category = "economic"
	Cultural  Category = "cultural"
	Tech      Category = "tech"
)

func CategoryFromInt16(v int16) Category {
	switch v {
	case 0:
		return Political
	case 1:
		return Social
	case 2:
		return Economic
	case 3:
		return Cultural
	case 4:
		return Tech
	}
	return ""
}

func (c Category) ToInt16() int16 {
	switch c {
	case Political:
		return 0
	case Social:
		return 1
	case Economic:
		return 2
	case Cultural:
		return 3
	case Tech:
		return 4
	}
	return -1
}

func TypeFromInt16(v int16) Type {
	switch v {
	case 0:
		return Arc
	case 1:
		return Event
	case 2:
		return Person
	case 3:
		return Object
	}
	return ""
}

func (t Type) ToInt16() int16 {
	switch t {
	case Arc:
		return 0
	case Event:
		return 1
	case Person:
		return 2
	case Object:
		return 3
	}
	return -1
}

func RecordStatusFromInt16(v int16) RecordStatus {
	switch v {
	case 0:
		return Removed
	case 1:
		return Draft
	case 2:
		return PendingReview
	case 3:
		return Reviewed
	}
	return ""
}

func (r RecordStatus) ToInt16() int16 {
	switch r {
	case Removed:
		return 0
	case Draft:
		return 1
	case PendingReview:
		return 2
	case Reviewed:
		return 3
	}
	return -1
}

func (s RecordService) Create(context context.Context, command createRecordCommandBody) (recordResponseBody, error) {
	response, err := s.recordRepository.Create(context, RecordAggregate{
		Record: model.Record{
			Title:        command.Title,
			Description:  command.Description,
			Location:     &command.Location,
			Significance: &command.Significance,
			URL:          command.Url,
			StartDate:    common.ToTime(command.StartDate),
			EndDate:      common.ToTime(command.EndDate),
			Type:         command.Type.ToInt16(),
			Status:       command.RecordStatus.ToInt16(),
		},
		Impacts: mapCreateImpacts(command.Impacts),
	})
	if err != nil {
		return recordResponseBody{}, err
	}
	return toResponse(response), nil
}

func mapUpdateImpacts(impactCommands []updateImpactCommandBody) []struct{ model.Impact } {
	impacts := make([]struct{ model.Impact }, len(impactCommands))
	for i, command := range impactCommands {
		impacts[i].Impact = model.Impact{
			ID:          command.ID,
			Description: command.Description,
			Value:       int16(command.Value),
			Category:    command.Category.ToInt16(),
			RecordID:    command.RecordId,
		}
	}
	return impacts
}

func mapCreateImpacts(impactCommands []createImpactCommandBody) []struct{ model.Impact } {
	impacts := make([]struct{ model.Impact }, len(impactCommands))
	for i, command := range impactCommands {
		impacts[i].Impact = model.Impact{
			ID:          uuid.New(),
			Description: command.Description,
			Value:       int16(command.Value),
			Category:    command.Category.ToInt16(),
		}
	}
	return impacts
}

func (s RecordService) GetById(id uuid.UUID) (recordResponseBody, error) {
	record, err := s.recordRepository.GetById(id)
	if err != nil {
		return recordResponseBody{}, err
	}
	return toResponse(record), nil
}

func (s RecordService) Update(c context.Context, id uuid.UUID, command updateRecordCommandBody) error {
	if id != command.ID {
		return errors.New("id mismatch")
	}
	return s.recordRepository.Update(c, RecordAggregate{
		Record: model.Record{
			ID:           command.ID,
			Title:        command.Title,
			Description:  command.Description,
			Location:     &command.Location,
			Significance: &command.Significance,
			URL:          command.Url,
			StartDate:    common.ToTime(command.StartDate),
			EndDate:      common.ToTime(command.EndDate),
			Type:         command.Type.ToInt16(),
			Status:       command.RecordStatus.ToInt16(),
		},
		Impacts: mapUpdateImpacts(command.Impacts),
	})
}

func (s RecordService) GetPaged(c context.Context, page, pageSize int) ([]recordResponseBody, int, error) {
	records, total, err := s.recordRepository.GetPaged(c, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	var response []recordResponseBody
	for _, record := range records {
		response = append(response, toResponse(record))
	}
	return response, total, nil
}

func (s RecordService) Delete(c context.Context, id uuid.UUID) error {
	return s.recordRepository.Delete(c, id)
}
