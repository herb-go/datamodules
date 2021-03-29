package storagetestutil

import (
	"bytes"
	"sync"
	"testing"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata"
)

type namespaces struct {
	lock     sync.Mutex
	storages map[string]*testStorage
}

func (ns *namespaces) Start() error {
	return nil
}
func (ns *namespaces) Stop() error {
	return nil
}
func (ns *namespaces) ExecuteGet(c herbcache.Context, key []byte) ([]byte, error) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	var err error
	buf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(buf, c.Namespace())
	if err != nil {
		return nil, err
	}
	s, ok := ns.storages[string(buf.Bytes())]
	if !ok {
		return nil, herbdata.ErrNotFound
	}
	return s.get(c, c.Position().RootDirectory(), key)
}

func (ns *namespaces) ExecuteSetWithTTL(c herbcache.Context, key []byte, data []byte, ttl int64) error {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	var err error
	buf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(buf, c.Namespace())
	if err != nil {
		return err
	}
	s, ok := ns.storages[string(buf.Bytes())]
	if !ok {
		s = newTestStorage()
		ns.storages[string(buf.Bytes())] = s
	}
	return s.setWithTTL(c, c.Position().RootDirectory(), key, data)
}
func (ns *namespaces) ExecuteDelete(c herbcache.Context, key []byte) error {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	var err error
	buf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(buf, c.Namespace())
	if err != nil {
		return err
	}
	s, ok := ns.storages[string(buf.Bytes())]
	if !ok {
		return nil
	}
	return s.delete(c, c.Position().RootDirectory(), key)

}
func (ns *namespaces) ExecuteFlush(c herbcache.Context) error {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	var err error
	buf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(buf, c.Namespace())
	if err != nil {
		return err
	}
	s, ok := ns.storages[string(buf.Bytes())]
	if !ok {
		return nil
	}
	return s.flush(c, c.Position().RootDirectory())
}

type testStorage struct {
	data map[string][]byte
	sub  map[string]*testStorage
}

func (s *testStorage) get(c herbcache.Context, d *herbcache.Directory, key []byte) ([]byte, error) {
	var err error
	var buf = bytes.NewBuffer(nil)
	if d == nil {
		_, err = herbcache.WriteGroupedKey(buf, c.Group(), key)
		if err != nil {
			return nil, err
		}
		bs, ok := s.data[string(buf.Bytes())]
		if !ok {
			return nil, herbdata.ErrNotFound
		}
		return bs, nil
	}
	_, err = herbcache.WriteDirectory(buf, d)
	if err != nil {
		return nil, err
	}
	sub, ok := s.sub[string(buf.Bytes())]
	if !ok {
		return nil, herbdata.ErrNotFound
	}
	return sub.get(c, d.Next, key)
}
func (s *testStorage) setWithTTL(c herbcache.Context, d *herbcache.Directory, key []byte, data []byte) error {
	var err error
	var buf = bytes.NewBuffer(nil)
	if d == nil {
		_, err = herbcache.WriteGroupedKey(buf, c.Group(), key)
		if err != nil {
			return err
		}
		s.data[string(buf.Bytes())] = data
		return nil
	}
	_, err = herbcache.WriteDirectory(buf, d)
	if err != nil {
		return err
	}
	sub, ok := s.sub[string(buf.Bytes())]
	if !ok {
		sub = newTestStorage()
		s.sub[string(buf.Bytes())] = sub
	}
	return sub.setWithTTL(c, d.Next, key, data)
}
func (s *testStorage) delete(c herbcache.Context, d *herbcache.Directory, key []byte) error {
	var err error
	var buf = bytes.NewBuffer(nil)
	if d == nil {
		_, err = herbcache.WriteGroupedKey(buf, c.Group(), key)
		if err != nil {
			return err
		}
		key := string(buf.Bytes())
		delete(s.data, key)
		return nil
	}
	_, err = herbcache.WriteDirectory(buf, d)
	if err != nil {
		return err
	}
	sub, ok := s.sub[string(buf.Bytes())]
	if !ok {
		return nil
	}
	return sub.delete(c, d.Next, key)
}
func (s *testStorage) flush(c herbcache.Context, d *herbcache.Directory) error {
	var err error
	var buf = bytes.NewBuffer(nil)
	_, err = herbcache.WriteDirectory(buf, d)
	if err != nil {
		return err
	}
	key := string(buf.Bytes())
	if d.Next == nil {
		delete(s.sub, key)
		return nil
	}
	sub, ok := s.sub[key]
	if !ok {
		return nil
	}
	return sub.flush(c, d.Next)
}

func newTestStorage() *testStorage {
	return &testStorage{
		data: map[string][]byte{},
		sub:  map[string]*testStorage{},
	}
}

func newStorage() *herbcache.Storage {
	s := herbcache.NewStorage()
	s.Engine = &namespaces{
		storages: map[string]*testStorage{},
	}
	return s
}

func factory() *herbcache.Storage {
	return newStorage()
}

func TestTool(t *testing.T) {
	TestNotFlushable(factory, func(*herbcache.Storage) {}, func(v ...interface{}) { t.Fatal(v...) })
	TestFlushable(factory, func(*herbcache.Storage) {}, func(v ...interface{}) { t.Fatal(v...) })
}
