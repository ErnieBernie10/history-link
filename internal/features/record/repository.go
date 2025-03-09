package record

import (
	"context"
	"database/sql"
	"fmt"
	"historylink/.gen/historylink/public/model"
	. "historylink/.gen/historylink/public/table"
	"log/slog"
	"reflect"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/google/uuid"
)

func NewRepository(db *sql.DB, logger *slog.Logger) IRecordRepository {
	return RecordRepository{
		db:     db,
		logger: logger,
	}
}

type IRecordRepository interface {
	GetById(uuid.UUID) (RecordAggregate, error)
	Create(c context.Context, command RecordAggregate) (RecordAggregate, error)
	Update(c context.Context, command RecordAggregate) error
	Delete(c context.Context, id uuid.UUID) error
	GetPaged(c context.Context, limit int, offset int) ([]RecordAggregate, int, error)
}
type RecordRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

type RecordAggregate struct {
	model.Record
	History model.RecordHistory

	Impacts []struct {
		model.Impact
	}
}

func (r RecordRepository) GetById(id uuid.UUID) (RecordAggregate, error) {
	stmt := SELECT(
		Record.AllColumns,
		Impact.AllColumns,
		RecordHistory.AllColumns,
	).FROM(
		Record.
			LEFT_JOIN(Impact, Impact.RecordID.EQ(Record.ID)).
			LEFT_JOIN(
				RecordHistory,
				RecordHistory.RecordID.EQ(Record.ID).
					AND(RecordHistory.UpdatedAt.IN(
						SELECT(MAX(RecordHistory.UpdatedAt)).
							FROM(RecordHistory).
							WHERE(RecordHistory.RecordID.EQ(Record.ID)),
					)),
			),
	).WHERE(
		Record.ID.EQ(UUID(id)),
	)

	var dest RecordAggregate
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return dest, fmt.Errorf("error getting record: %w", err)
	}

	fmt.Printf("%+v\n", dest)
	return dest, nil
}

func (r RecordRepository) Create(c context.Context, command RecordAggregate) (RecordAggregate, error) {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return RecordAggregate{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var result RecordAggregate

	recordStmt := Record.INSERT(Record.MutableColumns).
		MODEL(command.Record).
		RETURNING(Record.AllColumns)

	if err = recordStmt.Query(tx, &result.Record); err != nil {
		return RecordAggregate{}, fmt.Errorf("error creating record: %w", err)
	}

	for i := range command.Impacts {
		command.Impacts[i].RecordID = result.ID
	}

	if len(command.Impacts) > 0 {
		impactStmt := Impact.INSERT(Impact.MutableColumns).
			MODELS(command.Impacts).
			RETURNING(Impact.AllColumns)

		var impacts []struct {
			model.Impact
		}
		if err = impactStmt.Query(tx, &impacts); err != nil {
			return RecordAggregate{}, fmt.Errorf("error creating impacts: %w", err)
		}

		result.Impacts = impacts
	}

	if err = tx.Commit(); err != nil {
		return RecordAggregate{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return result, nil
}

func EqualsRecord(a, b model.Record) bool {
	return a.Title == b.Title &&
		a.Description == b.Description &&
		reflect.DeepEqual(a.Location, b.Location) &&
		reflect.DeepEqual(a.Significance, b.Significance) &&
		a.URL == b.URL &&
		a.StartDate.Equal(*b.StartDate) &&
		a.EndDate.Equal(*b.EndDate) &&
		a.Type == b.Type &&
		a.Status == b.Status
}

func EqualsImpact(a, b model.Impact) bool {
	return a.Description == b.Description &&
		a.Value == b.Value &&
		a.Category == b.Category
}

func (r RecordRepository) Update(c context.Context, command RecordAggregate) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Handle impacts - first get existing impacts
	existingImpactsStmt := SELECT(Impact.AllColumns).
		FROM(Impact).
		WHERE(Impact.RecordID.EQ(UUID(command.ID)))

	var existingImpacts []struct {
		model.Impact
	}
	err = existingImpactsStmt.Query(tx, &existingImpacts)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error getting existing impacts: %w", err)
	}

	// Create maps for existing and new impacts
	existingImpactsMap := make(map[uuid.UUID]struct{ model.Impact })
	for _, impact := range existingImpacts {
		existingImpactsMap[impact.Impact.ID] = impact
	}

	newImpactsMap := make(map[uuid.UUID]struct{ model.Impact })
	for _, impact := range command.Impacts {
		if impact.Impact.ID != uuid.Nil {
			newImpactsMap[impact.Impact.ID] = impact
		}
	}

	// Process impacts in three groups:
	// 1. Impacts to update (exist in both maps)
	// 2. Impacts to delete (exist in existing but not in new)
	// 3. Impacts to insert (exist in new but not in existing)

	// 1. Update existing impacts that have changed
	for id, newImpact := range newImpactsMap {
		if existingImpact, exists := existingImpactsMap[id]; exists {
			// Only update if something changed
			if !EqualsImpact(newImpact.Impact, existingImpact.Impact) {
				updateStmt := Impact.UPDATE(Impact.Description, Impact.Value, Impact.Category).
					MODEL(newImpact.Impact).
					WHERE(Impact.ID.EQ(UUID(id)))

				_, err = updateStmt.Exec(tx)
				if err != nil {
					return fmt.Errorf("error updating impact: %w", err)
				}
			}

			// Remove from existingImpactsMap to track what's been processed
			delete(existingImpactsMap, id)
		}
	}

	// 2. Delete impacts that no longer exist
	for id := range existingImpactsMap {
		deleteStmt := Impact.DELETE().WHERE(Impact.ID.EQ(UUID(id)))
		_, err = deleteStmt.Exec(tx)
		if err != nil {
			return fmt.Errorf("error deleting impact: %w", err)
		}
	}

	// 3. Insert new impacts
	for _, impact := range command.Impacts {
		if impact.Impact.ID == uuid.Nil {
			impact.Impact.RecordID = command.ID

			insertStmt := Impact.INSERT(Impact.Description, Impact.Value, Impact.Category, Impact.RecordID).
				MODEL(impact.Impact)

			_, err = insertStmt.Exec(tx)
			if err != nil {
				return fmt.Errorf("error inserting impact: %w", err)
			}
		}
	}

	stmt := SELECT(Record.AllColumns).
		FROM(Record).
		WHERE(Record.ID.EQ(UUID(command.ID)))

	var existingRecord struct {
		model.Record
	}
	err = stmt.Query(tx, &existingRecord)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("error getting existing record: %w", err)
	}

	if EqualsRecord(command.Record, existingRecord.Record) {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing transaction: %w", err)
		}
		return nil
	}

	// Update the record
	recordStmt := Record.UPDATE(Record.Title, Record.Description, Record.Location, Record.Significance, Record.URL, Record.StartDate, Record.EndDate, Record.Type, Record.Status).
		MODEL(command.Record).
		WHERE(Record.ID.EQ(UUID(command.ID)))

	_, err = recordStmt.Exec(tx)
	if err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r RecordRepository) Delete(c context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete the record
	recordStmt := Record.DELETE().WHERE(Record.ID.EQ(UUID(id)))
	_, err = recordStmt.Exec(tx)
	if err != nil {
		return fmt.Errorf("error deleting record: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

type Count struct {
	C int
}

func (r RecordRepository) GetPaged(c context.Context, limit int, offset int) ([]RecordAggregate, int, error) {
	var total Count
	stmt := SELECT(COUNT(Record.ID).AS("count.c")).FROM(Record)

	err := stmt.Query(r.db, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %w", err)
	}

	stmt = SELECT(
		Record.AllColumns,
		Impact.AllColumns,
		RecordHistory.AllColumns,
	).FROM(
		Record.
			LEFT_JOIN(Impact, Impact.RecordID.EQ(Record.ID)).
			LEFT_JOIN(RecordHistory, RecordHistory.RecordID.EQ(Record.ID)),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	var dest []RecordAggregate
	err = stmt.Query(r.db, &dest)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting records: %w", err)
	}
	return dest, total.C, nil
}
