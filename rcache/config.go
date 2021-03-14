package rcache

import (
	"github.com/herb-go/herbdata/kvdb"
)

type Config struct {
	Store       *kvdb.Config
	VersionTTL  int64
	Irrevocable bool
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
	e.VersionTTL = c.VersionTTL
	if e.VersionTTL == 0 {
		e.VersionTTL = DefaultVersionTTL
	}
	cache.CopyFrom(
		New().
			WithEngine(e).
			WithIrrevocable(c.Irrevocable),
	)
	return nil
}
