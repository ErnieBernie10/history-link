package common

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrLinkNotFound      = errors.New("link not found")
	ErrImpactNotFound    = errors.New("impact not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
	ErrLinkToItself      = errors.New("cannot link record to itself")
)
