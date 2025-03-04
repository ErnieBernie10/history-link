package record

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func NewRecordResources(conn *sql.DB, logger *slog.Logger) RecordResources {
	return RecordResources{
		logger:        logger,
		RecordService: NewRecordService(NewRepository(conn, logger), logger),
	}
}

type RecordResources struct {
	RecordService IRecordService
	logger        *slog.Logger
}

func (rs RecordResources) create(c context.Context, input *struct {
	Body createRecordCommandBody
}) (*struct {
	Body recordResponseBody
}, error) {
	response, err := rs.RecordService.Create(c, input.Body)
	if err != nil {
		rs.logger.Error(err.Error())
		return nil, err
	}

	return &struct {
		Body recordResponseBody
	}{
		Body: response,
	}, nil
}

func (rs RecordResources) getById(c context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*struct {
	Body recordResponseBody
}, error) {
	record, err := rs.RecordService.GetById(input.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, huma.Error404NotFound(fmt.Sprintf("Record with id %v not found", input.ID.String()))
		default:
			return nil, err
		}
	}

	return &struct {
		Body recordResponseBody
	}{
		Body: record,
	}, nil
}

func (rs RecordResources) update(c context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body updateRecordCommandBody
}) (*struct{}, error) {
	err := rs.RecordService.Update(c, input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	return &struct{}{}, nil
}

func (rs RecordResources) getPaged(c context.Context, input *struct {
	Page     int `query:"page" minimum:"1" default:"1"`
	PageSize int `query:"pageSize" minimum:"1" default:"10"`
}) (*struct {
	Body pagedResponse[recordResponseBody]
}, error) {
	records, total, err := rs.RecordService.GetPaged(c, input.Page, input.PageSize)
	if err != nil {
		return nil, err
	}

	if records == nil {
		records = []recordResponseBody{}
	}

	return &struct {
		Body pagedResponse[recordResponseBody]
	}{
		Body: pagedResponse[recordResponseBody]{
			Page:    input.Page,
			Size:    input.PageSize,
			Total:   total,
			Records: records,
		},
	}, nil
}

func (rs RecordResources) delete(c context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*struct{}, error) {
	err := rs.RecordService.Delete(c, input.ID)
	if err != nil {
		return nil, err
	}

	return &struct{}{}, nil
}

func (rs RecordResources) MountRoutes(s huma.API) {
	huma.Register(s, huma.Operation{
		OperationID:   "get-record-by-id",
		Method:        http.MethodGet,
		Path:          "/records/{id}",
		DefaultStatus: http.StatusOK,
	}, rs.getById)
	huma.Register(s, huma.Operation{
		OperationID: "create-record",
		Method:      http.MethodPost,
		Path:        "/records/",
	}, rs.create)
	huma.Register(s, huma.Operation{
		OperationID: "update-record",
		Method:      http.MethodPut,
		Path:        "/records/{id}",
	}, rs.update)
	huma.Register(s, huma.Operation{
		OperationID:   "get-records",
		Method:        http.MethodGet,
		Path:          "/records/",
		DefaultStatus: http.StatusOK,
	}, rs.getPaged)
	huma.Register(s, huma.Operation{
		OperationID: "delete-record",
		Method:      http.MethodDelete,
		Path:        "/records/{id}",
	}, rs.delete)
}
