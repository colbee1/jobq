package service

import "errors"

var (
	ErrInvalidAction = errors.New("jqs: invalid action")
	ErrInvalidStatus = errors.New("jqs: invalid status")
)
