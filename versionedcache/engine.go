package versionedcache

import (
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/datautil"
)

type Engine struct {
	VersionGenerator func() (string, error)
	VersionTTL       int64
	VersionStore     herbdata.SetterGetterServer
	Store            herbdata.CacheServer
}

func (e *Engine) LoadRawVersion(key []byte) ([]byte, error) {
	v, err := e.VersionStore.Get(key)
	if err == nil {
		return v, nil
	}
	if err == herbdata.ErrNotFound {
		return []byte{}, nil
	}
	return nil, err
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
	var vererr error
	var err error
	if e.VersionStore != nil {
		vererr = e.VersionStore.Stop()
	}
	err = e.Store.Stop()
	if vererr != nil {
		return vererr
	}
	if err != nil {
		return err
	}
	return nil
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
