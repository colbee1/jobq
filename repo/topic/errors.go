package topic

import "errors"

var (
	ErrTopicNotFound          = errors.New("topic repo: Topic not found")
	ErrStatsCollectorDisabled = errors.New("topic repo: Stats collector is disabled")
)
