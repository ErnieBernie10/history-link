package link

import (
	"historylink/.gen/historylink/public/model"

	"github.com/google/uuid"
)

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

func mapLinkResponseBody(m model.Link, index int) linkResponseBody {
	return linkResponseBody{
		ID:       m.ID,
		RecordID: m.RecordId2,
		Strength: m.Strength,
	}
}
