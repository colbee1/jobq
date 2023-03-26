package badger3

import (
	"github.com/dgraph-io/badger/v3"
)

type Adapter struct {
	db       *badger.DB
	jobIdSeq *badger.Sequence
}

type Options struct {
	SyncWrite bool
	DropAll   bool // Only for testing
}
