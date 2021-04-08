package cachetool

import (
	"time"

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

var DefaultEncoding = msgpackencoding.Encoding

type Tool struct {
	TTL       time.Duration
	Cache     herbdata.Cache
	LockerMap *waitingmap.LockerMap
	Encoding  *dataencoding.Encoding
}

func NewTool() *Tool {
	return &Tool{}
}
