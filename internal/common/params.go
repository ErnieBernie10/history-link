package common

import (
	"errors"

	"github.com/go-fuego/fuego"
)

func ExtractPageParams(c fuego.ContextNoBody) (int, int, error) {
	page, err1 := c.QueryParamIntErr("page")
	pageSize, err2 := c.QueryParamIntErr("pageSize")
	err := errors.Join(err1, err2)
	if err != nil {
		return 0, 0, fuego.BadRequestError{
			Err: err,
		}
	}

	errorItems := []fuego.ErrorItem{}
	if page < 0 {
		errorItems = append(errorItems, fuego.ErrorItem{Name: "page", Reason: "page must be greater than or equal to 0"})
	}
	if pageSize < 1 {
		errorItems = append(errorItems, fuego.ErrorItem{Name: "pageSize", Reason: "pageSize must be greater than 1"})
	}

	if len(errorItems) > 0 {
		return 0, 0, fuego.BadRequestError{
			Err:    errors.New("validation error with page params"),
			Errors: errorItems,
		}
	}

	return page, pageSize, nil
}
