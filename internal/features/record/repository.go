package record

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"historylink/internal/db"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

func NewRepository(db *sql.DB) IRecordRepository {
	return RecordRepository{
		db: db,
	}
}

type IRecordRepository interface {
	GetById(uuid.UUID) (Record, error)
	Create(c context.Context, command Record) (Record, error)
	Update(c context.Context, command Record) error
	GetPaged(c context.Context, limit int, offset int) ([]Record, error)
}
type RecordRepository struct {
	db *sql.DB
}

type Impact struct {
	ID          uuid.UUID `db:"id" fieldtag:"pk"`
	RecordID    uuid.UUID `db:"record_id"`
	Description string    `db:"description"`
	Value       int       `db:"value"`
	Category    int       `db:"category"`
}

type Record struct {
	ID           uuid.UUID      `db:"id" fieldtag:"pk"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Location     sql.NullString `db:"location"`
	Significance sql.NullString `db:"significance"`
	Url          string         `db:"url"`
	StartDate    sql.NullTime   `db:"start_date"`
	EndDate      sql.NullTime   `db:"end_date"`
	Type         int16          `db:"type"`
	RecordStatus int16          `db:"status"`
	Impacts      []Impact       `db:"-"`
}

func EqualsImpact(a, b Impact) bool {
	return a.ID == b.ID &&
		a.RecordID == b.RecordID &&
		a.Description == b.Description &&
		a.Value == b.Value &&
		a.Category == b.Category
}

func EqualsRecord(a, b Record) bool {
	return a.ID == b.ID &&
		a.Title == b.Title &&
		a.Description == b.Description &&
		a.Location.String == b.Location.String &&
		a.Significance.String == b.Significance.String &&
		a.Url == b.Url &&
		a.StartDate == b.StartDate &&
		a.EndDate == b.EndDate &&
		a.Type == b.Type &&
		a.RecordStatus == b.RecordStatus
}

var recordStruct = sqlbuilder.NewStruct(new(Record)).For(sqlbuilder.PostgreSQL)
var impactStruct = sqlbuilder.NewStruct(new(Impact)).For(sqlbuilder.PostgreSQL)

func (r RecordRepository) GetById(u uuid.UUID) (Record, error) {
	sb := recordStruct.SelectFrom("record")
	sb.Select("*")
	sb.Where(sb.EQ("id", u.String()))

	q, args := sb.Build()

	row := r.db.QueryRow(q, args...)

	var record Record
	if err := db.ScanStructRow(row, &record, recordStruct); err != nil {
		return Record{}, err
	}

	var err error
	if record.Impacts, err = r.getImpactsForRecord(record.ID); err != nil {
		return Record{}, err
	}

	return record, nil
}

func (r RecordRepository) getImpactsForRecord(recordId uuid.UUID) ([]Impact, error) {
	sb := impactStruct.SelectFrom("impact")
	sb.Select("*")
	sb.Where(sb.EQ("record_id", recordId.String()))

	q, args := sb.Build()

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return []Impact{}, err
	}
	defer rows.Close()

	var impacts []Impact
	if err := db.ScanStructRows(rows, &impacts, impactStruct); err != nil {
		return []Impact{}, err
	}

	return impacts, nil
}

func (r RecordRepository) Create(c context.Context, command Record) (Record, error) {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return Record{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the record first
	record, err := r.insertRecord(tx, command)
	if err != nil {
		return Record{}, fmt.Errorf("failed to insert record: %w", err)
	}

	// If there are impacts, insert them
	if len(command.Impacts) > 0 {
		impacts, err := r.insertImpacts(tx, record.ID, command.Impacts)
		if err != nil {
			return Record{}, fmt.Errorf("failed to insert impacts: %w", err)
		}
		record.Impacts = impacts
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return Record{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return record, nil
}

// insertRecord inserts a new record into the database
func (r RecordRepository) insertRecord(tx *sql.Tx, command Record) (Record, error) {
	sb := recordStruct.InsertInto("record")
	sb.Cols("title", "description", "location", "significance", "url", "start_date", "end_date", "type", "status")
	sb.Values(command.Title, command.Description, command.Location, command.Significance, command.Url, command.StartDate, command.EndDate, command.Type, command.RecordStatus)
	sb.Returning("*")

	q, args := sb.Build()

	row := tx.QueryRow(q, args...)

	var record Record
	if err := db.ScanStructRow(row, &record, recordStruct); err != nil {
		return Record{}, err
	}

	return record, nil
}

// insertImpacts inserts impacts for a record into the database
func (r RecordRepository) insertImpacts(tx *sql.Tx, recordID uuid.UUID, impactCommands []Impact) ([]Impact, error) {
	sb := impactStruct.InsertInto("impact")
	sb.Cols("record_id", "description", "value", "category")

	for _, impact := range impactCommands {
		sb.Values(recordID, impact.Description, impact.Value, impact.Category)
	}
	sb.Returning("*")

	q, args := sb.Build()

	rows, err := tx.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var impacts []Impact
	if err := db.ScanStructRows(rows, &impacts, impactStruct); err != nil {
		return nil, err
	}

	return impacts, nil
}

func (r RecordRepository) updateImpacts(tx *sql.Tx, recordID uuid.UUID, impactCommands []Impact) error {
	if len(impactCommands) == 0 {
		return nil
	}

	// First, get all the current impacts for this record
	sb := impactStruct.SelectFrom("impact")
	sb.Select("*")
	sb.Where(sb.EQ("record_id", recordID.String()))

	q, args := sb.Build()

	rows, err := tx.Query(q, args...)
	if err != nil {
		return fmt.Errorf("failed to fetch current impacts: %w", err)
	}
	defer rows.Close()

	var currentImpacts []Impact
	if err := db.ScanStructRows(rows, &currentImpacts, impactStruct); err != nil {
		return fmt.Errorf("failed to scan current impacts: %w", err)
	}

	// Create a map of current impacts by ID for easy lookup
	currentImpactsMap := make(map[string]Impact)
	for _, impact := range currentImpacts {
		currentImpactsMap[impact.ID.String()] = impact
	}

	// Update only impacts that have changed
	for _, impactCommand := range impactCommands {
		// Check if this impact exists in the current impacts
		currentImpact, exists := currentImpactsMap[impactCommand.ID.String()]

		// Skip update if nothing has changed
		if exists &&
			EqualsImpact(currentImpact, impactCommand) {
			continue
		}

		// If we got here, we need to update this impact
		sb := impactStruct.Update("impact", impactCommand)
		sb.Where(sb.EQ("id", impactCommand.ID.String()))
		q, args := sb.Build()

		if _, err := tx.Exec(q, args...); err != nil {
			return fmt.Errorf("failed to update impact %s: %w", impactCommand.ID, err)
		}
	}

	return nil
}

// Update updates a record in the database
func (r RecordRepository) Update(c context.Context, command Record) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update impacts
	if len(command.Impacts) > 0 {
		if err := r.updateImpacts(tx, command.ID, command.Impacts); err != nil {
			return fmt.Errorf("failed to update impacts: %w", err)
		}
	}

	if r, err := r.GetById(command.ID); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("record with id %s not found: %w", command.ID, err)
	} else if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	} else if EqualsRecord(r, command) {
		return nil
	}

	// Update the record
	sb := recordStruct.Update("record", command)
	sb.Where(sb.EQ("id", command.ID.String()))

	q, args := sb.Build()

	_, err = tx.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r RecordRepository) GetPaged(c context.Context, limit int, offset int) ([]Record, error) {
	sb := recordStruct.SelectFrom("record r")
	sb.Select("*")
	sb.Join("impact i", "r.id = i.record_id")
	sb.Limit(limit)
	sb.Offset(offset)

	q, args := sb.Build()

	fmt.Println(q, args)

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recordsMap := map[uuid.UUID]Record{}
	for rows.Next() {
		var record Record
		if err := rows.Scan(recordStruct.Addr(&record)...); err != nil {
			return nil, err
		}

		// Check if the record already exists
		if exists, ok := recordsMap[record.ID]; ok {
			exists.Impacts = append(exists.Impacts, record.Impacts...)
		} else {
			recordsMap[record.ID] = record
		}
	}

	records := make([]Record, 0, len(recordsMap))
	for _, record := range recordsMap {
		records = append(records, record)
	}

	return records, nil
}
