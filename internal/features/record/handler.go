package record

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func NewRecordResources(conn *sql.DB) RecordResources {
	return RecordResources{
		RecordService: NewRecordService(NewRepository(conn)),
	}
}

type RecordResources struct {
	RecordService IRecordService
}

func (rs RecordResources) create(c context.Context, input *struct {
	Body createRecordCommandBody
}) (*struct {
	Body recordResponseBody
}, error) {
	response, err := rs.RecordService.Create(c, input.Body)
	if err != nil {
		return nil, err
	}

	return &struct {
		Body recordResponseBody
	}{
		Body: toResponse(response),
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
		Body: toResponse(record),
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
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}) (*struct {
	Body pagedResponse[recordResponseBody]
}, error) {
	records, err := rs.RecordService.GetPaged(c, input.Page, input.PageSize)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, huma.Error404NotFound("No records found")
		default:
			return nil, err
		}
	}

	return &struct {
		Body pagedResponse[recordResponseBody]
	}{
		Body: toPagedResponse(records, input.Page, input.PageSize),
	}, nil
}

func (rs RecordResources) MountRoutes(s huma.API) {
	//fuego.Get(s, "/record/", rs.getPaged, option.QueryInt("page", "page of records", param.Default(1), param.Nullable()), option.QueryInt("pageSize", "size of page", param.Default(10), param.Nullable()))
	huma.Register(s, huma.Operation{
		OperationID:   "get-record-by-id",
		Method:        http.MethodGet,
		Path:          "/record/{id}",
		DefaultStatus: http.StatusOK,
	}, rs.getById)
	huma.Register(s, huma.Operation{
		OperationID: "create-record",
		Method:      http.MethodPost,
		Path:        "/record",
	}, rs.create)
	huma.Register(s, huma.Operation{
		OperationID: "update-record",
		Method:      http.MethodPut,
		Path:        "/record/{id}",
	}, rs.update)
	huma.Register(s, huma.Operation{
		OperationID:   "get-records",
		Method:        http.MethodGet,
		Path:          "/record/",
		DefaultStatus: http.StatusOK,
	}, rs.getPaged)
}
