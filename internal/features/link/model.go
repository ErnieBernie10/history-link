package link

import "github.com/google/uuid"

type createLinkCommandBody struct {
	RecordID uuid.UUID `json:"recordId"`
	Strength int16     `json:"strength"`
}

type updateLinkCommandBody struct {
	Strength int16 `json:"strength"`
}

type linkResponseBody struct {
	ID       uuid.UUID `json:"id"`
	RecordID uuid.UUID `json:"recordId"`
	Strength int16     `json:"strength"`
}
