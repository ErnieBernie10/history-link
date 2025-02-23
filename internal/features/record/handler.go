package record

import (
	"database/sql"
	"errors"
	"fmt"
	"historylink/internal/common"
	data "historylink/internal/db"
	"historylink/internal/util"
	"net/http"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func NewRecordResources(db *data.Queries) RecordResources {
	return RecordResources{
		RecordService: NewRecordService(db),
	}
}

type RecordResources struct {
	RecordService IRecordService
}

func (rs RecordResources) create(c fuego.ContextWithBody[createRecordCommandBody]) (recordResponse, error) {
	req, err := c.Body()
	if err != nil {
		return recordResponse{}, err
	}

	start, err1 := time.Parse("200601021504", req.StartDate)
	end, err2 := time.Parse("200601021504", req.EndDate)
	err = errors.Join(err1, err2)

	if err != nil {
		common.ToHttpErrors(err)
	}

	response, err := rs.RecordService.Create(createRecordCommand{
		Title:        req.Title,
		Description:  req.Description,
		Location:     req.Location,
		Significance: req.Significance,
		Url:          req.Url,
		StartDate:    start,
		EndDate:      end,
		Status:       req.Status,
		Type:         req.Type,
	})
	if err != nil {
		return recordResponse{}, err
	}

	return toResponse(response), nil
}

func (rs RecordResources) getById(c fuego.ContextNoBody) (recordResponse, error) {
	id, err := util.PathParamUuid(c, "id")
	if err != nil {
		return recordResponse{}, err
	}

	record, err := rs.RecordService.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return recordResponse{}, fuego.NotFoundError{
				Err: fmt.Errorf("Record with id %v not found", id.String()),
			}
		default:
			return recordResponse{}, fuego.InternalServerError{Err: err}
		}
	}

	return toResponse(record), nil
}
func (rs RecordResources) MountRoutes(s *fuego.Server) {
	fuego.Get(s, "/{id}", rs.getById)
	fuego.Post(s, "/", rs.create, option.DefaultStatusCode(http.StatusCreated))
}
