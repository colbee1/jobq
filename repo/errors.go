package repo

import "errors"

var (
	ErrInvalidTransaction = errors.New("job repo: invalid transaction")
	ErrJobNotFound        = errors.New("job repo: not found")
	ErrInvalidRetryCount  = errors.New("job repo: invalid retry count")
)
