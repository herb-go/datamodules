package rcache

import (
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/datautil"
)

type Engine struct {
	VersionGenerator func() (string, error)
	Store            herbdata.CacheServer
	VersionStore     herbdata.StoreServer
}

func (e *Engine) Start() error {
	if e.VersionStore != nil {
		err := e.VersionStore.Start()
		if err != nil {
			return err
		}
	}
	return e.Store.Start()
}
func (e *Engine) Stop() error {
	if e.VersionStore != nil {
		err := e.VersionStore.Stop()
		if err != nil {
			return err
		}
	}
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
	}
}
