package link

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"historylink/.gen/historylink/public/model"
	. "historylink/.gen/historylink/public/table"
	"historylink/internal/common"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type ILinkRepository interface {
	Create(c context.Context, command model.Link) (model.Link, error)
	GetById(id uuid.UUID) (model.Link, error)
	GetByRecordId(c context.Context, recordId uuid.UUID) ([]model.Link, error)
	GetByRecordIds(c context.Context, recordId uuid.UUID, recordId2 uuid.UUID) (model.Link, error)
	Update(c context.Context, id uuid.UUID, command model.Link) error
	Delete(c context.Context, id uuid.UUID) error
}

type LinkRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) ILinkRepository {
	return LinkRepository{
		db:     db,
		logger: logger,
	}
}

func (r LinkRepository) Create(c context.Context, command model.Link) (model.Link, error) {
	selectStmt := Record.SELECT(Record.ID).
		FROM(Record).
		WHERE(Record.ID.EQ(UUID(command.RecordID)).OR(Record.ID.EQ(UUID(command.RecordId2))))

	var records []model.Record
	if err := selectStmt.Query(r.db, &records); err != nil {
		if err == sql.ErrNoRows {
			return model.Link{}, common.ErrRecordNotFound
		}
	}
	if len(records) != 2 {
		return model.Link{}, common.ErrRecordNotFound
	}

	stmt := Link.INSERT(Link.MutableColumns).
		MODEL(command).
		RETURNING(Link.AllColumns)

	var dest model.Link
	if err := stmt.Query(r.db, &dest); err != nil {
		return model.Link{}, fmt.Errorf("failed to create link: %w", err)
	}

	return dest, nil
}

func (r LinkRepository) GetById(id uuid.UUID) (model.Link, error) {
	stmt := SELECT(Link.AllColumns).
		FROM(Link).
		WHERE(Link.ID.EQ(UUID(id)))

	var dest model.Link
	if err := stmt.Query(r.db, &dest); err != nil {
		return model.Link{}, fmt.Errorf("failed to get link by id: %w", err)
	}

	return dest, nil
}

func (r LinkRepository) GetByRecordId(c context.Context, recordId uuid.UUID) ([]model.Link, error) {
	stmt := SELECT(Link.AllColumns).
		FROM(Link).
		WHERE(Link.RecordID.EQ(UUID(recordId)).
			OR(Link.RecordId2.EQ(UUID(recordId))))

	var dest []model.Link
	if err := stmt.Query(r.db, &dest); err != nil {
		return nil, fmt.Errorf("failed to get link by record id: %w", err)
	}

	return dest, nil
}

func (r LinkRepository) GetByRecordIds(c context.Context, recordId uuid.UUID, recordId2 uuid.UUID) (model.Link, error) {
	stmt := SELECT(Link.AllColumns).
		FROM(Link).
		WHERE(Link.RecordID.EQ(UUID(recordId)).
			AND(Link.RecordId2.EQ(UUID(recordId2))).
			OR(Link.RecordID.EQ(UUID(recordId2)).
				AND(Link.RecordId2.EQ(UUID(recordId)))))

	fmt.Println(stmt.DebugSql())

	var dest model.Link
	if err := stmt.Query(r.db, &dest); err != nil {
		if err == sql.ErrNoRows {
			return model.Link{}, fmt.Errorf("link not found: %w", err)
		}
		return model.Link{}, fmt.Errorf("failed to get link by record ids: %w", err)
	}

	return dest, nil
}

func (r LinkRepository) Update(c context.Context, id uuid.UUID, command model.Link) error {
	stmt := Link.UPDATE(Link.Strength).
		SET(command).
		WHERE(Link.ID.EQ(UUID(id)))

	if _, err := stmt.Exec(r.db); err != nil {
		return fmt.Errorf("failed to update link: %w", err)
	}

	return nil
}

func (r LinkRepository) Delete(c context.Context, id uuid.UUID) error {
	stmt := Link.DELETE().
		WHERE(Link.ID.EQ(UUID(id)))

	if _, err := stmt.Exec(r.db); err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}

	return nil
}
