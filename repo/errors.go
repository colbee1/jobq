package repo

import "errors"

var (
	ErrInvalidTransaction = errors.New("job repo: invalid transaction")
	ErrJobNotFound        = errors.New("job repo: not found")
)
