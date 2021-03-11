package rcache

import (
	"github.com/herb-go/herbdata/kvdb"
)

type Config struct {
	Store        *kvdb.Config
	TTL          int64
	Prefix       string
	VersionStore *kvdb.Config
	Irrevocable  bool
}

func (c *Config) ApplyTo(cache *Cache) error {
	var err error
	e := NewEngine()
	db := kvdb.New()
	err = c.Store.ApplyTo(db)
	if err != nil {
		return err
	}
	e.Store = db
	if c.VersionStore != nil {
		vdb := kvdb.New()
		err = c.VersionStore.ApplyTo(vdb)
		e.VersionStore = vdb
	}
	cache.CopyFrom(
		New().
			WithEngine(e).
			WithTTL(c.TTL).
			WithPath([]byte(c.Prefix)).
			WithIrrevocable(c.Irrevocable),
	)
	return nil
}
