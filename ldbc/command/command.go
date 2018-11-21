package command

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Command interface {
	Parse(args []string) error
	Run(db *leveldb.DB) error
}
