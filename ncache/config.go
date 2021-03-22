package ncache

import (
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/herbdata/kvdb/commonkvdb"
)

type StorageConfig struct {
	Cache        *kvdb.Config
	VersionStore *kvdb.Config
	VersionTTL   int64
}

func (c *StorageConfig) ApplyTo(storage *Storage) error {
	var err error
	var versiondb *kvdb.Database
	db := kvdb.New()
	err = c.Cache.ApplyTo(db)
	if err != nil {
		return err
	}
	if !db.Features().SupportAll(kvdb.FeatureTTLStore) {
		return kvdb.ErrFeatureNotSupported
	}
	storage.Cache = db
	if c.VersionStore != nil && c.VersionStore.Driver != "" {
		versiondb = kvdb.New()
		err = c.VersionStore.ApplyTo(versiondb)
		if err != nil {
			return err
		}
		storage.VersionStore = versiondb
	} else {
		features := db.Features()
		if features.SupportAll(kvdb.FeatureStore) && !features.SupportAny(kvdb.FeatureUnstable) {
			storage.VersionStore = db
		} else if features.SupportAll(kvdb.FeatureEmbedded | kvdb.FeatureNonpersistent) {
			versiondb = kvdb.New()
			err = commonkvdb.NewInMemory().ApplyTo(versiondb)
			if err != nil {
				return err
			}
			storage.VersionStore = versiondb
		}
	}
	storage.VersionTTL = c.VersionTTL
	return nil
}

func NewStorageConfig() *StorageConfig {
	return &StorageConfig{}
}
