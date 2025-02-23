package subject

import (
	data "berniestack/internal/db"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func NewSubjectService(db *data.Queries) ISubjectService {
	return SubjectService{
		db: db,
	}
}

type ISubjectService interface {
	Create(command createSubjectCommand)
	GetById(id uuid.UUID) (data.Subject, error)
}

type SubjectService struct {
	db *data.Queries
}

type createSubjectCommand struct {
	Title       string
	Summary     string
	SubjectType string
	Url         string
	Weight      int
	FromDate    time.Time
	UntilDate   time.Time
}

func (s SubjectService) Create(command createSubjectCommand) {
	s.db.CreateSubject(context.Background(), data.CreateSubjectParams{
		ID:          uuid.NewString(),
		Title:       command.Title,
		Summary:     command.Summary,
		Weight:      sql.NullInt64{Valid: true, Int64: int64(command.Weight)},
		SubjectType: sql.NullString{Valid: true, String: command.SubjectType},
		Url:         command.Url,
		FromDate:    command.FromDate,
		UntilDate:   command.UntilDate,
	})
}

func (s SubjectService) GetById(id uuid.UUID) (data.Subject, error) {
	return s.db.GetSubject(context.Background(), id.String())
}
