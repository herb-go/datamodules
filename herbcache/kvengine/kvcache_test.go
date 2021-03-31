package kvengine

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/datamodules/herbcache/storagetestutil"
	_ "github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/kvdb"
	_ "github.com/herb-go/herbdata/kvdb/commonkvdb"
)

var factory = func() *herbcache.Storage {
	s := herbcache.NewStorage()
	config := &StorageConfig{
		Cache: &kvdb.Config{
			Driver: "freecache",
			Config: func(v interface{}) error {
				return json.Unmarshal([]byte(`{"Size":50000}`), v)
			},
		},
		VersionTTL: 3600,
		VersionStore: &kvdb.Config{
			Driver: "inmemory",
		},
	}

	err := config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	return s
}

func novcachefactory() *herbcache.Storage {
	s := factory()
	s.Engine.(*Engine).VersionTTL = 0
	return s
}
func TestKVCache(t *testing.T) {
	storagetestutil.TestNotFlushable(factory, func(*herbcache.Storage) {}, func(v ...interface{}) { t.Fatal(v...) })
	storagetestutil.TestFlushable(factory, func(*herbcache.Storage) {}, func(v ...interface{}) { t.Fatal(v...) })
	storagetestutil.TestFlushable(factory, func(*herbcache.Storage) {}, func(v ...interface{}) { t.Fatal(v...) })

}
