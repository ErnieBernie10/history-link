package record

import (
	"context"
	"errors"
	"historylink/.gen/historylink/public/model"
	"time"

	"github.com/google/uuid"
)

func NewRecordService(recordRepository IRecordRepository) IRecordService {
	return RecordService{
		recordRepository: recordRepository,
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

func (s RecordService) Create(context context.Context, command createRecordCommandBody) (recordResponseBody, error) {
	startDate, err := time.Parse("200601021504", command.StartDate)
	if err != nil {
		return recordResponseBody{}, err
	}
	endDate, err := time.Parse("200601021504", command.EndDate)
	if err != nil {
		return recordResponseBody{}, err
	}
	response, err := s.recordRepository.Create(context, RecordAggregate{
		Record: model.Record{
			Title:        command.Title,
			Description:  command.Description,
			Location:     &command.Location,
			Significance: &command.Significance,
			URL:          command.Url,
			StartDate:    &startDate,
			EndDate:      &endDate,
			Type:         int16(command.Type),
			Status:       int16(command.RecordStatus),
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
			Category:    int16(command.Category),
			RecordID:    &command.RecordId,
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
			Category:    int16(command.Category),
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
	startDate, err := time.Parse("200601021504", command.StartDate)
	if err != nil {
		return err
	}
	endDate, err := time.Parse("200601021504", command.EndDate)
	if err != nil {
		return err
	}
	return s.recordRepository.Update(c, RecordAggregate{
		Record: model.Record{
			ID:           command.ID,
			Title:        command.Title,
			Description:  command.Description,
			Location:     &command.Location,
			Significance: &command.Significance,
			URL:          command.Url,
			StartDate:    &startDate,
			EndDate:      &endDate,
			Type:         int16(command.Type),
			Status:       int16(command.RecordStatus),
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
