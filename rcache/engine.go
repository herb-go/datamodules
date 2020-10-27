package rcache

import (
	"time"

	"github.com/herb-go/herbdata/datautil"
)

type Engine struct {
	versionGenerator func() (string, error)
	versionstore     VersionStore
	driver           Driver
}

func (e *Engine) NewCache(path []byte, irrevocable bool) *Cache {
	return &Cache{
		Path:        path,
		irrevocable: irrevocable,
		engine:      e,
	}
}

var DefaultVersionbGenerator = func() (string, error) {
	v, err := datautil.Encode(uint64(time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func NewEngine() *Engine {
	return &Engine{
		versionGenerator: DefaultVersionbGenerator,
	}
}
