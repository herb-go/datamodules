package cachepreset_test

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/datamodules/herbcache/cachepreset"
	"github.com/herb-go/datamodules/herbcache/kvengine"
	"github.com/herb-go/herbdata"
	_ "github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/dataencoding/msgpackencoding"
	"github.com/herb-go/herbdata/kvdb"
	_ "github.com/herb-go/herbdata/kvdb/commonkvdb"
)

func testcache() *herbcache.Cache {
	cache := herbcache.New()
	storage := herbcache.NewStorage()
	config := &kvengine.StorageConfig{
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
	err := config.ApplyTo(storage)
	if err != nil {
		panic(err)
	}
	return cache.OverrideStorage(storage).OverrideFlushable(true)
}
func TestOperation(t *testing.T) {
	cache := testcache()
	enc := msgpackencoding.Encoding
	preset, err := cachepreset.New(cachepreset.Cache(cache), cachepreset.Encoding(enc), cachepreset.TTL(100)).Apply()
	if err != nil {
		t.Fatal()
	}
	if preset.Cache() != cache || preset.Encoding() != enc {
		t.Fatal()
	}
	_, err = preset.GetS("test")
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = preset.SetWithTTLS("test", []byte("testvalue"), 100)
	if err != nil {
		t.Fatal(err)
	}
	data, err := preset.GetS("test")
	if err != nil || string(data) != "testvalue" {
		t.Fatal(err)
	}
	err = preset.DeleteS("test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = preset.GetS("test")
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	var result string
	err = preset.StoreS("test2", "testvalue2", 100)
	if err != nil {
		t.Fatal(err)
	}
	err = preset.LoadS("test2", &result)
	if err != nil || result != "testvalue2" {
		t.Fatal(err)
	}

	pc, err := preset.Concat(cachepreset.PrefixCache("prefix")).Apply()
	if err != nil {
		t.Fatal(err)
	}
	err = pc.StoreS("test2", "prefixtestvalue2", 100)
	if err != nil {
		t.Fatal(err)
	}
	err = preset.LoadS("test2", &result)
	if err != nil || result != "testvalue2" {
		t.Fatal(err)
	}
	err = pc.LoadS("test2", &result)
	if err != nil || result != "prefixtestvalue2" {
		t.Fatal(err)
	}
	cc, err := preset.Concat(cachepreset.ChildCache("prefix")).Apply()
	if err != nil {
		t.Fatal(err)
	}
	err = cc.StoreS("test2", "childtestvalue2", 100)
	if err != nil {
		t.Fatal(err)
	}
	err = cc.LoadS("test2", &result)
	if err != nil || result != "childtestvalue2" {
		t.Fatal(err)
	}
	cc2, err := cc.Concat(cachepreset.ChildCache("prefix2")).Apply()
	if err != nil {
		t.Fatal(err)
	}
	err = cc2.StoreS("test2", "child2testvalue2", 100)
	if err != nil {
		t.Fatal(err)
	}
	err = cc2.LoadS("test2", &result)
	if err != nil || result != "child2testvalue2" {
		t.Fatal(err)
	}
	err = cc.Flush()
	if err != nil {
		t.Fatal(err)
	}
	err = preset.LoadS("test2", &result)
	if err != nil || result != "testvalue2" {
		t.Fatal(err)
	}
	err = pc.LoadS("test2", &result)
	if err != nil || result != "prefixtestvalue2" {
		t.Fatal(err)
	}
	err = cc.LoadS("test2", &result)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
	err = cc2.LoadS("test2", &result)
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
}
