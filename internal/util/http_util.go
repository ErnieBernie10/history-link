package util

import (
	"errors"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func PathParamUuid(c fuego.ContextWithBody[any], paramname string) (uuid.UUID, error) {
	param := c.PathParam(paramname)
	paramValue, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, fuego.BadRequestError{
			Err: errors.New("Invalid route param"),
			Errors: []fuego.ErrorItem{
				{Name: paramname, Reason: "Invalid parameter"},
			},
		}
	}
	return paramValue, nil
}
