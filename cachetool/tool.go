package cachetool

import (
	"time"

	"github.com/herb-go/herbdata/kvdb"

	"github.com/herb-go/herbdata/dataencoding/msgpackencoding"

	"github.com/herb-go/herbdata/dataencoding"
	"github.com/herb-go/misc/waitingmap"

	"github.com/herb-go/herbdata"
)

type Loader interface {
	Load(key []byte) error
}

type Saver interface {
	Save(key []byte) error
}

type LoaderSaver interface {
	Loader
	Saver
}

var DefaultEncoding *dataencoding.Encoding = msgpackencoding.Encoding

var DefaultTTL time.Duration = time.Hour

var DefaultCache herbdata.Cache = kvdb.Passthrough

type Tool struct {
	TTL       time.Duration
	Cache     herbdata.Cache
	LockerMap *waitingmap.LockerMap
	Encoding  *dataencoding.Encoding
}

func NewTool() *Tool {
	return &Tool{
		TTL:      DefaultTTL,
		Encoding: DefaultEncoding,
		Cache:    DefaultCache,
	}
}
