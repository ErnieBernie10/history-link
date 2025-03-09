//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type ImpactHistory struct {
	ID          uuid.UUID `sql:"primary_key"`
	ImpactID    *uuid.UUID
	RecordID    *uuid.UUID
	Description string
	Value       int16
	Category    int16
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
