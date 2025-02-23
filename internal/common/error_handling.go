package common

import (
	"errors"

	"github.com/go-fuego/fuego"
)

func ToHttpErrors(err error) []fuego.ErrorItem {
	var errorItems []fuego.ErrorItem

	// If it's a multi-error from errors.Join, extract all wrapped errors
	var joinedErrors interface{ Unwrap() []error }
	if errors.As(err, &joinedErrors) {
		for _, e := range joinedErrors.Unwrap() {
			errorItems = append(errorItems, fuego.ErrorItem{Reason: e.Error()})
		}
		return errorItems
	}

	// Otherwise, handle single wrapped errors
	for err != nil {
		errorItems = append(errorItems, fuego.ErrorItem{Reason: err.Error()})
		err = errors.Unwrap(err)
	}

	return errorItems
}
