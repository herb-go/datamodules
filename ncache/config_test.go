package ncache

import (
	"testing"

	"github.com/herb-go/herbdata/kvdb/commonkvdb"

	"github.com/herb-go/herbdata/kvdb"
)

type testDriver struct {
	kvdb.Nop
	Feature kvdb.Feature
}

//Features return supported features
func (d testDriver) Features() kvdb.Feature {
	return d.Feature
}

func newTestDriverConfig(feature kvdb.Feature) *Config {
	config := NewConfig()
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

func getCacheVersionStoreDirver(c *Cache) kvdb.Driver {
	s := c.engine.VersionStore.(*kvdb.Database)
	return s.Driver
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
	var config *Config
	var c *Cache
	config = newTestDriverConfig(0)
	c = New()
	err = config.ApplyTo(c)
	if err != kvdb.ErrFeatureNotSupported {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore)
	c = New()
	err = config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
	if c.engine.VersionStore != nil {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore | kvdb.FeatureStore)
	c = New()
	err = config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
	if c.engine.VersionStore == nil {
		t.Fatal()
	}
	if _, ok := getCacheVersionStoreDirver(c).(testDriver); !ok {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore | kvdb.FeatureStore | kvdb.FeatureUnstable)
	c = New()
	err = config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
	if c.engine.VersionStore != nil {
		t.Fatal()
	}
	config = newTestDriverConfig(kvdb.FeatureTTLStore | kvdb.FeatureEmbedded | kvdb.FeatureNonpersistent)
	c = New()
	err = config.ApplyTo(c)
	if err != nil {
		panic(err)
	}
	if c.engine.VersionStore == nil {
		t.Fatal()
	}
	if _, ok := getCacheVersionStoreDirver(c).(*commonkvdb.InMemory); !ok {
		t.Fatal()
	}
}
