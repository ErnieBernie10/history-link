package subject

import (
	data "berniestack/internal/db"
	"berniestack/internal/util"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-fuego/fuego"
)

func NewSubjectResources(db *data.Queries) SubjectResources {
	return SubjectResources{
		SubjectService: NewSubjectService(db),
	}
}

type SubjectResources struct {
	SubjectService ISubjectService
}

type createSubjectRequest struct {
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	SubjectType string `json:"subjectType"`
	Url         string `json:"url"`
	Weight      int    `json:"weight"`
	FromDate    string `json:"fromDate" validate:"required,datetime=200601021504"`
	UntilDate   string `json:"untilDate" validate:"required,datetime=200601021504"`
}

func (rs SubjectResources) create(c fuego.ContextWithBody[createSubjectRequest]) (any, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}
	fromDate, err1 := time.Parse("200601021504", req.FromDate)
	untilDate, err2 := time.Parse("200601021504", req.UntilDate)
	err = errors.Join(err1, err2)
	if err != nil {
		return nil, err
	}

	rs.SubjectService.Create(createSubjectCommand{
		Url:         req.Url,
		SubjectType: req.SubjectType,
		Title:       req.Title,
		Summary:     req.Summary,
		Weight:      req.Weight,
		FromDate:    fromDate,
		UntilDate:   untilDate})
	return nil, nil
}

func (rs SubjectResources) getById(c fuego.ContextNoBody) (data.Subject, error) {
	id, err := util.PathParamUuid(c, "id")
	if err != nil {
		return data.Subject{}, err
	}

	subject, err := rs.SubjectService.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return subject, fuego.NotFoundError{
				Err: fmt.Errorf("Subject with id %v not found", id.String()),
			}
		default:
			return subject, fuego.InternalServerError{Err: err}
		}
	}

	return subject, nil
}
func (rs SubjectResources) MountRoutes(s *fuego.Server) {
	fuego.Get(s, "/{id}", rs.getById)
	fuego.Post(s, "/", rs.create)
}
