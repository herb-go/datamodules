package versionedcache

import (
	"github.com/herb-go/herbdata/kvdb"
)

type Config struct {
	Store        *kvdb.Config
	VersionStore *kvdb.Config
	VersionTTL   int64
	Revocable    bool
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
		versiondb := kvdb.New()
		err = c.VersionStore.ApplyTo(versiondb)
		if err != nil {
			return err
		}
		e.VersionStore = versiondb
	}
	e.VersionTTL = c.VersionTTL
	cache.CopyFrom(
		New().
			WithEngine(e).
			WithRevocable(c.Revocable),
	)
	return nil
}
