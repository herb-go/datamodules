package storagetestutil

import (
	"bytes"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata"
)

func TestNotFlushable(creator func() herbcache.Storage, closer func(herbcache.Storage), fatal func(...interface{})) {
	var data []byte
	var err error
	s := creator()
	err = s.Start()
	if err != nil {
		fatal(err)
	}
	defer func() {
		err = s.Stop()
		if err != nil {
			fatal(err)
		}
	}()
	defer closer(s)
	c := herbcache.New().OverrideStorage(s)
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey2"))
	if err != nil {
		panic(err)
	}
	err = c.SetWithTTL([]byte("testkey"), nil, 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, nil) || err != nil {
		fatal(data, err)
	}
	c = c.Migrate([]byte("bs"))
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	cg := c.OverrideGroup([]byte("group"))
	data, err = cg.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = cg.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = cg.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = cg.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = cg.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	c = c.SubCache([]byte("sub1"))
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	c = c.SubCache([]byte("sub2"))
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).SetWithTTL([]byte("test2test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).OverrideGroup([]byte("test2")).SetWithTTL([]byte("test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).SubCache([]byte("test2")).SetWithTTL([]byte("test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
}

func TestFlushable(creator func() herbcache.Storage, closer func(herbcache.Storage), fatal func(...interface{})) {
	var data []byte
	var err error
	s := creator()
	err = s.Start()
	if err != nil {
		fatal(err)
	}
	defer func() {
		err = s.Stop()
		if err != nil {
			fatal(err)
		}
	}()
	defer closer(s)
	c := herbcache.New().OverrideStorage(s).OverrideFlushable(true)
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey2"))
	if err != nil {
		panic(err)
	}
	err = c.SetWithTTL([]byte("testkey"), nil, 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, nil) || err != nil {
		fatal(data, err)
	}
	c = c.Migrate([]byte("bs"))
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	cg := c.OverrideGroup([]byte("group"))
	data, err = cg.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = cg.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = cg.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = cg.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = cg.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	c = c.SubCache([]byte("sub1"))
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	c = c.SubCache([]byte("sub2"))
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SetWithTTL([]byte("testkey"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if !bytes.Equal(data, []byte(data)) || err != nil {
		fatal(data, err)
	}
	err = c.Delete([]byte("testkey"))
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testkey"))
	if data != nil || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).SetWithTTL([]byte("test2test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).OverrideGroup([]byte("test2")).SetWithTTL([]byte("test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).SubCache([]byte("test2")).SetWithTTL([]byte("test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	c = c.Migrate(nil)
	data, err = c.SubCache([]byte("test1")).Get([]byte("test2"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).SetWithTTL([]byte("test2"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1")).Get([]byte("test2"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).Flush()
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1")).Get([]byte("test2"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	data, err = c.SubCache([]byte("test1")).SubCache([]byte("test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).SetWithTTL([]byte("test2"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	err = c.SubCache([]byte("test1")).SubCache([]byte("test2")).SetWithTTL([]byte("test3"), []byte("data"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1")).SubCache([]byte("test2")).Get([]byte("test3"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	data, err = c.SubCache([]byte("test1")).Get([]byte("test2"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	err = c.SubCache([]byte("test1")).Flush()
	if err != nil {
		panic(err)
	}
	data, err = c.SubCache([]byte("test1")).SubCache([]byte("test2")).Get([]byte("test3"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}

	data, err = c.SubCache([]byte("test1")).Get([]byte("test2"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}

	err = c.SetWithTTL([]byte("testns"), []byte("testdata"), 3600)
	if err != nil {
		panic(err)
	}
	err = c.Migrate([]byte("n")).SetWithTTL([]byte("testns"), []byte("testdata"), 3600)
	if err != nil {
		panic(err)
	}
	err = c.OverrideGroup([]byte("n")).SetWithTTL([]byte("testns"), []byte("testdata"), 3600)
	if err != nil {
		panic(err)
	}
	data, err = c.Get([]byte("testns"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	data, err = c.Migrate([]byte("n")).Get([]byte("testns"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	data, err = c.OverrideGroup([]byte("n")).Get([]byte("testns"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	err = c.Flush()
	if err != nil {
		fatal(err)
	}
	data, err = c.Get([]byte("testns"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
	data, err = c.Migrate([]byte("n")).Get([]byte("testns"))
	if len(data) == 0 || err != nil {
		fatal(data, err)
	}
	data, err = c.OverrideGroup([]byte("n")).Get([]byte("testns"))
	if len(data) != 0 || err != herbdata.ErrNotFound {
		fatal(data, err)
	}
}
