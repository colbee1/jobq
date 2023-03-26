package jobq

import (
	"github.com/pkg/errors"
)

var (
	ErrTopicIsMissing   = errors.New("jobq: topic is missing")
	ErrTopicNotFound    = errors.New("jobq: topic not found")
	ErrJobNotFound      = errors.New("jobq: job not found")
	ErrInvalidJobStatus = errors.New("jobq: invalid job state")
)
