package record

import (
	"context"
	"database/sql"
	"historylink/.gen/historylink/public/model"
	. "historylink/.gen/historylink/public/table"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/google/uuid"
)

func NewRepository(db *sql.DB) IRecordRepository {
	return RecordRepository{
		db: db,
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
	db *sql.DB
}

type RecordAggregate struct {
	model.Record
	Impacts []struct {
		model.Impact
	}
}

func EqualsImpact(a, b struct{ model.Impact }) bool {
	return a.Impact.ID == b.Impact.ID &&
		a.Impact.Description == b.Impact.Description &&
		a.Impact.Value == b.Impact.Value &&
		a.Impact.Category == b.Impact.Category
}

func (r RecordRepository) GetById(id uuid.UUID) (RecordAggregate, error) {
	stmt := SELECT(
		Record.AllColumns,
		Impact.AllColumns,
	).FROM(
		Record.LEFT_JOIN(Impact, Impact.RecordID.EQ(Record.ID)),
	).WHERE(
		Record.ID.EQ(UUID(id)),
	)

	var dest RecordAggregate
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

func (r RecordRepository) Create(c context.Context, command RecordAggregate) (RecordAggregate, error) {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return RecordAggregate{}, err
	}
	defer tx.Rollback()

	// Create a result struct to hold the returned data
	var result RecordAggregate

	// Insert the record and capture the returned data
	recordStmt := Record.INSERT(Record.Title, Record.Description, Record.Location, Record.Significance, Record.URL, Record.StartDate, Record.EndDate, Record.Type, Record.Status).
		MODEL(command.Record).
		RETURNING(Record.AllColumns)

	err = recordStmt.Query(tx, &result.Record)
	if err != nil {
		return RecordAggregate{}, err
	}

	// Prepare impacts with the correct record ID if needed
	for i := range command.Impacts {
		if command.Impacts[i].RecordID == nil {
			command.Impacts[i].RecordID = &result.ID
		}
	}

	// Insert impacts if any
	if len(command.Impacts) > 0 {
		impactStmt := Impact.INSERT(Impact.Description, Impact.Value, Impact.Category, Impact.RecordID).
			MODELS(command.Impacts).
			RETURNING(Impact.AllColumns)

		var impacts []struct {
			model.Impact
		}
		err = impactStmt.Query(tx, &impacts)
		if err != nil {
			return RecordAggregate{}, err
		}
		result.Impacts = impacts
	}

	err = tx.Commit()
	if err != nil {
		return RecordAggregate{}, err
	}

	return result, nil
}

func (r RecordRepository) Update(c context.Context, command RecordAggregate) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the record
	recordStmt := Record.UPDATE(Record.Title, Record.Description, Record.Location, Record.Significance, Record.URL, Record.StartDate, Record.EndDate, Record.Type, Record.Status).
		MODEL(command.Record).
		WHERE(Record.ID.EQ(UUID(command.ID)))

	_, err = recordStmt.Exec(tx)
	if err != nil {
		return err
	}

	// Handle impacts - first get existing impacts
	existingImpactsStmt := SELECT(Impact.AllColumns).
		FROM(Impact).
		WHERE(Impact.RecordID.EQ(UUID(command.ID)))

	var existingImpacts []struct {
		model.Impact
	}
	err = existingImpactsStmt.Query(tx, &existingImpacts)
	if err != nil && err != sql.ErrNoRows {
		return err
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
			if !EqualsImpact(newImpact, existingImpact) {
				updateStmt := Impact.UPDATE(Impact.Description, Impact.Value, Impact.Category).
					MODEL(newImpact.Impact).
					WHERE(Impact.ID.EQ(UUID(id)))

				_, err = updateStmt.Exec(tx)
				if err != nil {
					return err
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
			return err
		}
	}

	// 3. Insert new impacts
	for _, impact := range command.Impacts {
		if impact.Impact.ID == uuid.Nil {
			impact.Impact.RecordID = &command.ID

			insertStmt := Impact.INSERT(Impact.Description, Impact.Value, Impact.Category, Impact.RecordID).
				MODEL(impact.Impact)

			_, err = insertStmt.Exec(tx)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r RecordRepository) Delete(c context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the record
	recordStmt := Record.DELETE().WHERE(Record.ID.EQ(UUID(id)))
	_, err = recordStmt.Exec(tx)
	if err != nil {
		return err
	}
	tx.Commit()

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
		return nil, 0, err
	}

	stmt = SELECT(
		Record.AllColumns,
		Impact.AllColumns,
	).FROM(
		Record.LEFT_JOIN(Impact, Impact.RecordID.EQ(Record.ID)),
	).LIMIT(int64(limit)).
		OFFSET(int64(offset))

	var dest []RecordAggregate
	err = stmt.Query(r.db, &dest)
	if err != nil {
		return nil, 0, err
	}
	return dest, total.C, nil
}
