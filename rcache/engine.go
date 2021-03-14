package rcache

import (
	"time"

	"github.com/herb-go/misc/waitingmap"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/datautil"
)

var DefaultVersionTTL = int64(36000)

type Engine struct {
	VersionGenerator func() (string, error)
	VersionTTL       int64
	Store            herbdata.CacheServer
	oncemap          *waitingmap.OnceMap
}

func (e *Engine) Start() error {
	return e.Store.Start()
}
func (e *Engine) Stop() error {
	return e.Store.Stop()
}

var DefaultVersionGenerator = func() (string, error) {
	v, err := datautil.Encode(uint64(time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func NewEngine() *Engine {
	return &Engine{
		VersionGenerator: DefaultVersionGenerator,
		oncemap:          waitingmap.NewOnceMap(),
	}
}
