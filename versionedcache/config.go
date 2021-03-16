package versionedcache

import (
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/herbdata/kvdb/commonkvdb"
)

type Config struct {
	Store        *kvdb.Config
	VersionStore *kvdb.Config
	VersionTTL   int64
	Revocable    bool
}

func (c *Config) ApplyTo(cache *Cache) error {
	var err error
	var versiondb *kvdb.Database
	e := NewEngine()
	db := kvdb.New()
	err = c.Store.ApplyTo(db)
	if err != nil {
		return err
	}
	e.Cache = db
	if c.VersionStore != nil {
		versiondb = kvdb.New()
		err = c.VersionStore.ApplyTo(versiondb)
		if err != nil {
			return err
		}
		e.VersionStore = versiondb
	} else {
		features := db.Features()
		if features.SupportAll(kvdb.FeatureStore) && !features.SupportAny(kvdb.FeatureUnstable) {
			e.VersionStore = db
		} else if features.SupportAll(kvdb.FeatureEmbedded | kvdb.FeatureNonpersistent) {
			versiondb = kvdb.New()
			err = commonkvdb.NewInMemory().ApplyTo(versiondb)
			if err != nil {
				return err
			}
			e.VersionStore = versiondb
		}
	}
	if c.Revocable && e.VersionStore == nil {
		return ErrNoVersionStore
	}
	e.VersionTTL = c.VersionTTL
	cache.CopyFrom(
		New().
			WithEngine(e).
			WithRevocable(c.Revocable),
	)
	return nil
}
