package jobq

import (
	"errors"
	"fmt"
)

var (
	ErrPackage = errors.New("jobq")
)

var (
	ErrSlugIsEmpty     = fmt.Errorf("%w: job slug is missing", ErrPackage)
	ErrSlugAlreadySet  = fmt.Errorf("%w: job slug is already set", ErrPackage)
	ErrTopicIsInvalid  = fmt.Errorf("%w: queue topic is invalid", ErrPackage)
	ErrTopicNotFound   = fmt.Errorf("%w: topic not found", ErrPackage)
	ErrJobNotFound     = fmt.Errorf("%w: job not found", ErrPackage)
	ErrInvalidJobState = fmt.Errorf("%w: invalid job state", ErrPackage)
)
