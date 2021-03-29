package kvcache

import (
	"bytes"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/datautil"
)

type Storage struct {
	VersionGenerator func() (string, error)
	VersionTTL       int64
	VersionStore     herbdata.SetterGetterServer
	Cache            herbdata.CacheServer
}

func (s *Storage) Start() error {
	if s.VersionStore != nil {
		err := s.VersionStore.Start()
		if err != nil {
			return err
		}
	}
	return s.Cache.Start()
}
func (s *Storage) Stop() error {
	var vererr error
	var err error
	if s.VersionStore != nil {
		vererr = s.VersionStore.Stop()
	}
	err = s.Cache.Stop()
	if vererr != nil {
		return vererr
	}
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) rawKey(c herbcache.Context, key []byte) []byte {
	var err error
	buf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(buf, c.Namespace())
	if err != nil {
		panic(err)
	}
	_, err = herbcache.WriteDirectories(buf, c.Position().RootDirectory())
	if err != nil {
		panic(err)
	}
	_, err = herbcache.WriteGroupedKey(buf, c.Group(), key)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
func (s *Storage) loadVersion(key []byte, cacheable bool) ([]byte, error) {
	if cacheable {
		return s.getCachedVersion(key)
	}
	return s.LoadRawVersion(key)
}
func (s *Storage) getCachedVersion(key []byte) ([]byte, error) {
	version, err := s.Cache.Get(key)
	if err == nil {
		return version, nil
	}
	if err != herbdata.ErrNotFound {
		return nil, err
	}
	version, err = s.LoadRawVersion(key)
	if err != nil {
		return nil, err
	}
	err = s.Cache.SetWithTTL(key, version, s.VersionTTL)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (s *Storage) LoadRawVersion(key []byte) ([]byte, error) {
	v, err := s.VersionStore.Get(key)
	if err == nil {
		return v, nil
	}
	if err == herbdata.ErrNotFound {
		return []byte{}, nil
	}
	return nil, err
}
func (s *Storage) getRawkeyAndVersion(c herbcache.Context, key []byte) ([]byte, []byte, error) {
	var err error
	versionbuf := bytes.NewBuffer(nil)
	keybuf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(keybuf, c.Namespace())
	if err != nil {
		panic(err)
	}
	cacheable := s.VersionTTL > 0 && s.VersionStore != nil
	d := c.Position().RootDirectory()
	for d != nil {
		_, err = herbcache.WriteDirectory(keybuf, d)
		if err != nil {
			return nil, nil, err
		}
		currentkey := keybuf.Bytes()
		v, err := s.loadVersion(currentkey, cacheable)
		if err != nil {
			return nil, nil, err
		}
		err = datautil.PackTo(versionbuf, nil, v)
		if err != nil {
			return nil, nil, err
		}
		keybuf = bytes.NewBuffer(currentkey)
		d = d.Next
	}
	_, err = herbcache.WriteGroupedKey(keybuf, c.Group(), key)
	if err != nil {
		return nil, nil, err
	}
	return keybuf.Bytes(), versionbuf.Bytes(), nil
}
func (s *Storage) setVersion(c herbcache.Context, version []byte) error {
	var err error
	cacheable := s.VersionTTL > 0 && s.VersionStore != nil
	keybuf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(keybuf, c.Namespace())
	if err != nil {
		return err
	}
	_, err = herbcache.WriteDirectories(keybuf, c.Position().RootDirectory())
	if err != nil {
		return err
	}
	key := keybuf.Bytes()
	err = s.VersionStore.Set(key, version)
	if err != nil {
		return err
	}
	if cacheable {
		return s.Cache.Delete(key)
	}
	return nil
}
func (s *Storage) ExecuteGet(c herbcache.Context, key []byte) ([]byte, error) {
	var data []byte
	var version []byte
	var err error
	var e *enity
	var rawkey []byte
	flushable := c.Flushable()
	if flushable {
		rawkey, version, err = s.getRawkeyAndVersion(c, key)
		if err != nil {
			return nil, err
		}
	} else {
		rawkey = s.rawKey(c, key)
	}
	data, err = s.Cache.Get(rawkey)
	if err != nil {
		return nil, err
	}
	e, err = loadEnity(data, flushable, version)
	if err != nil {
		if err == ErrEnityTypecodeNotMatch || err == ErrEnityVersionNotMatch {
			return nil, herbdata.ErrNotFound
		}
		return nil, err
	}
	return e.data, nil
}
func (s *Storage) ExecuteSetWithTTL(c herbcache.Context, key []byte, data []byte, ttl int64) error {
	var version []byte
	var err error
	var e *enity
	var rawkey []byte
	flushable := c.Flushable()
	if flushable {
		rawkey, version, err = s.getRawkeyAndVersion(c, key)
		if err != nil {
			return err
		}
	} else {
		rawkey = s.rawKey(c, key)
	}
	e = createEnity(flushable, version, data)
	buf := bytes.NewBuffer(nil)
	err = e.SaveTo(buf)
	if err != nil {
		return err
	}
	return s.Cache.SetWithTTL(rawkey, buf.Bytes(), ttl)
}
func (s *Storage) ExecuteDelete(c herbcache.Context, key []byte) error {
	return s.Cache.Delete(s.rawKey(c, key))
}
func (s *Storage) ExecuteFlush(c herbcache.Context) error {
	if !c.Flushable() {
		return herbdata.ErrNotFlushable
	}
	if s.VersionStore == nil {
		return ErrNoVersionStore
	}
	v, err := s.VersionGenerator()
	if err != nil {
		return err
	}
	return s.setVersion(c, []byte(v))

}

func New() *Storage {
	return &Storage{
		VersionGenerator: DefaultVersionGenerator,
	}
}
