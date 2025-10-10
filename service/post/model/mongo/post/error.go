package model

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidObjectId = errors.New("invalid objectId")
)
