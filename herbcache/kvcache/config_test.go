package kvcache

import (
	"testing"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/herbdata/kvdb/commonkvdb"
)

type testDriver struct {
	kvdb.Nop
	Feature kvdb.Feature
}

//Features return supported features
func (d testDriver) Features() kvdb.Feature {
	return d.Feature
}

func newTestDriverConfig(feature kvdb.Feature) *StorageConfig {
	config := NewStorageConfig()
	config.Cache = &kvdb.Config{
		Driver: "vcachetestdriver",
		Config: func(v interface{}) error {
			d := v.(*testDriver)
			d.Feature = feature
			return nil
		},
	}
	return config
}

func getCacheVersionStoreDirver(s *Engine) kvdb.Driver {
	return s.VersionStore.(*kvdb.Database).Driver
}
func init() {
	kvdb.Register("vcachetestdriver", func(loader func(v interface{}) error) (kvdb.Driver, error) {
		driver := testDriver{}
		err := loader(&driver)
		if err != nil {
			return nil, err
		}
		return driver, nil
	})
}

func TestConfig(t *testing.T) {
	var err error
	var config *StorageConfig
	var s *herbcache.Storage
	config = newTestDriverConfig(0)
	s = herbcache.NewStorage()
	err = config.ApplyTo(s)
	if err != kvdb.ErrFeatureNotSupported {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore)
	s = herbcache.NewStorage()
	err = config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	if s.Engine.(*Engine).VersionStore != nil {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore | kvdb.FeatureStore)
	s = herbcache.NewStorage()
	err = config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	if s.Engine.(*Engine).VersionStore == nil {
		t.Fatal()
	}
	if _, ok := getCacheVersionStoreDirver(s.Engine.(*Engine)).(testDriver); !ok {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore | kvdb.FeatureStore | kvdb.FeatureUnstable)
	s = herbcache.NewStorage()
	err = config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	if s.Engine.(*Engine).VersionStore != nil {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore | kvdb.FeatureEmbedded | kvdb.FeatureNonpersistent)
	s = herbcache.NewStorage()
	err = config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	if s.Engine.(*Engine).VersionStore == nil {
		t.Fatal()
	}
	if _, ok := getCacheVersionStoreDirver(s.Engine.(*Engine)).(*commonkvdb.InMemory); !ok {
		t.Fatal()
	}
}
