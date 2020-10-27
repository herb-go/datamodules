package rcache

import (
	"github.com/herb-go/herbdata"
)

type Driver interface {
	herbdata.Cache
	herbdata.Closer
	VersionStore
}
