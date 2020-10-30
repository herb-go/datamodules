package rcache

import (
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/datautil"
)

type Engine struct {
	VersionGenerator func() (string, error)
	Store            herbdata.Cache
	VersionStore     herbdata.Store
	lockerMap        *datautil.LockerMap
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
		VersionGenerator: DefaultVersionbGenerator,
		lockerMap:        datautil.NewLockerMap(),
	}
}
