package link

import (
	"context"
	"log/slog"

	"historylink/.gen/historylink/public/model"
	"historylink/internal/common"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ILinkService interface {
	Create(c context.Context, command createLinkCommandBody, targetRecordId uuid.UUID) (linkResponseBody, error)
	GetById(id uuid.UUID) (linkResponseBody, error)
	GetByRecordId(c context.Context, recordId uuid.UUID) ([]linkResponseBody, error)
	Delete(c context.Context, id uuid.UUID) error
}

type LinkService struct {
	linkRepository ILinkRepository
	logger         *slog.Logger
}

func NewLinkService(linkRepository ILinkRepository, logger *slog.Logger) ILinkService {
	return LinkService{
		linkRepository: linkRepository,
		logger:         logger,
	}
}

func (s LinkService) Create(c context.Context, command createLinkCommandBody, targetRecordId uuid.UUID) (linkResponseBody, error) {
	if command.RecordID == targetRecordId {
		return linkResponseBody{}, common.ErrLinkToItself
	}
	if _, err := s.linkRepository.GetByRecordIds(c, targetRecordId, command.RecordID); err == nil {
		return linkResponseBody{}, common.ErrLinkAlreadyExists
	}
	res, err := s.linkRepository.Create(c, model.Link{
		RecordID:  targetRecordId,
		RecordId2: command.RecordID,
		Strength:  command.Strength,
	})
	if err != nil {
		return linkResponseBody{}, err
	}

	return linkResponseBody{
		ID:       res.ID,
		RecordID: res.RecordId2,
		Strength: res.Strength,
	}, nil
}

func (s LinkService) GetById(id uuid.UUID) (linkResponseBody, error) {
	panic("implement me")
}

func (s LinkService) GetByRecordId(c context.Context, recordId uuid.UUID) ([]linkResponseBody, error) {
	links, err := s.linkRepository.GetByRecordId(c, recordId)
	if err != nil {
		return nil, err
	}

	return lo.Map(links, mapLinkResponseBody), nil
}

func (s LinkService) Delete(c context.Context, id uuid.UUID) error {
	return s.linkRepository.Delete(c, id)
}
