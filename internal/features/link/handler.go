package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"historylink/internal/common"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func NewLinkResources(conn *sql.DB, logger *slog.Logger) LinkResources {
	return LinkResources{
		logger:      logger,
		LinkService: NewLinkService(NewRepository(conn, logger), logger),
	}
}

type LinkResources struct {
	LinkService ILinkService
	logger      *slog.Logger
}

func (rs LinkResources) create(c context.Context, input *struct {
	ID   uuid.UUID `path:"record_id"`
	Body createLinkCommandBody
}) (*struct {
	Body linkResponseBody
}, error) {
	response, err := rs.LinkService.Create(c, input.Body, input.ID)
	if err != nil {
		switch err {
		case common.ErrLinkToItself:
			return nil, huma.Error400BadRequest(err.Error())
		case common.ErrLinkAlreadyExists:
			return nil, huma.Error409Conflict(err.Error())
		case common.ErrRecordNotFound:
			return nil, huma.Error404NotFound(err.Error())
		}
		rs.logger.Error(err.Error())
		return nil, err
	}

	return &struct {
		Body linkResponseBody
	}{
		Body: response,
	}, nil
}

func (rs LinkResources) getById(c context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*struct {
	Body linkResponseBody
}, error) {
	link, err := rs.LinkService.GetById(input.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, huma.Error404NotFound(fmt.Sprintf("Link with id %v not found", input.ID.String()))
		default:
			return nil, err
		}
	}

	return &struct {
		Body linkResponseBody
	}{
		Body: link,
	}, nil
}

func (rs LinkResources) getByRecordId(c context.Context, input *struct {
	RecordId uuid.UUID `path:"record_id"`
}) (*struct {
	Body []linkResponseBody
}, error) {
	links, err := rs.LinkService.GetByRecordId(c, input.RecordId)
	if err != nil {
		return nil, err
	}

	return &struct {
		Body []linkResponseBody
	}{
		Body: links,
	}, nil
}

func (rs LinkResources) delete(c context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*struct{}, error) {
	err := rs.LinkService.Delete(c, input.ID)
	if err != nil {
		return nil, err
	}

	return &struct{}{}, nil
}

func (rs LinkResources) MountRoutes(s huma.API) {
	huma.Register(s, huma.Operation{
		OperationID: "get-links-by-record-id",
		Method:      http.MethodGet,
		Path:        "/records/{record_id}/links",
	}, rs.getByRecordId)
	huma.Register(s, huma.Operation{
		OperationID: "create-link",
		Method:      http.MethodPost,
		Path:        "/records/{record_id}/links",
		Responses: map[string]*huma.Response{
			"409": {
				Description: "Link already exists",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: s.OpenAPI().Components.Schemas.Schema(reflect.TypeOf(huma.ErrorModel{}), true, ""),
					},
				},
			},
			"404": {
				Description: "Record not found",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: s.OpenAPI().Components.Schemas.Schema(reflect.TypeOf(huma.ErrorModel{}), true, ""),
					},
				},
			},
			"400": {
				Description: "Link to itself",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: s.OpenAPI().Components.Schemas.Schema(reflect.TypeOf(huma.ErrorModel{}), true, ""),
					},
				},
			},
		},
	}, rs.create)
	huma.Register(s, huma.Operation{
		OperationID: "delete-link",
		Method:      http.MethodDelete,
		Path:        "/links/{id}",
	}, rs.delete)
}
