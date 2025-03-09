//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var ImpactHistory = newImpactHistoryTable("public", "impact_history", "")

type impactHistoryTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnString
	ImpactID    postgres.ColumnString
	RecordID    postgres.ColumnString
	Description postgres.ColumnString
	Value       postgres.ColumnInteger
	Category    postgres.ColumnInteger
	CreatedAt   postgres.ColumnTimestamp
	UpdatedAt   postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ImpactHistoryTable struct {
	impactHistoryTable

	EXCLUDED impactHistoryTable
}

// AS creates new ImpactHistoryTable with assigned alias
func (a ImpactHistoryTable) AS(alias string) *ImpactHistoryTable {
	return newImpactHistoryTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ImpactHistoryTable with assigned schema name
func (a ImpactHistoryTable) FromSchema(schemaName string) *ImpactHistoryTable {
	return newImpactHistoryTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ImpactHistoryTable with assigned table prefix
func (a ImpactHistoryTable) WithPrefix(prefix string) *ImpactHistoryTable {
	return newImpactHistoryTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ImpactHistoryTable with assigned table suffix
func (a ImpactHistoryTable) WithSuffix(suffix string) *ImpactHistoryTable {
	return newImpactHistoryTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newImpactHistoryTable(schemaName, tableName, alias string) *ImpactHistoryTable {
	return &ImpactHistoryTable{
		impactHistoryTable: newImpactHistoryTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newImpactHistoryTableImpl("", "excluded", ""),
	}
}

func newImpactHistoryTableImpl(schemaName, tableName, alias string) impactHistoryTable {
	var (
		IDColumn          = postgres.StringColumn("id")
		ImpactIDColumn    = postgres.StringColumn("impact_id")
		RecordIDColumn    = postgres.StringColumn("record_id")
		DescriptionColumn = postgres.StringColumn("description")
		ValueColumn       = postgres.IntegerColumn("value")
		CategoryColumn    = postgres.IntegerColumn("category")
		CreatedAtColumn   = postgres.TimestampColumn("created_at")
		UpdatedAtColumn   = postgres.TimestampColumn("updated_at")
		allColumns        = postgres.ColumnList{IDColumn, ImpactIDColumn, RecordIDColumn, DescriptionColumn, ValueColumn, CategoryColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns    = postgres.ColumnList{ImpactIDColumn, RecordIDColumn, DescriptionColumn, ValueColumn, CategoryColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return impactHistoryTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		ImpactID:    ImpactIDColumn,
		RecordID:    RecordIDColumn,
		Description: DescriptionColumn,
		Value:       ValueColumn,
		Category:    CategoryColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
