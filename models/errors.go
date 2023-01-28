package models

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("item is not found")
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
)
