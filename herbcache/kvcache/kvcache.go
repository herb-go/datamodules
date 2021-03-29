package kvcache

import (
	"bytes"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/datautil"
)

type Engine struct {
	VersionGenerator func() (string, error)
	VersionTTL       int64
	VersionStore     herbdata.SetterGetterServer
	Cache            herbdata.CacheServer
}

func (e *Engine) Start() error {
	if e.VersionStore != nil {
		err := e.VersionStore.Start()
		if err != nil {
			return err
		}
	}
	return e.Cache.Start()
}
func (e *Engine) Stop() error {
	var vererr error
	var err error
	if e.VersionStore != nil {
		vererr = e.VersionStore.Stop()
	}
	err = e.Cache.Stop()
	if vererr != nil {
		return vererr
	}
	if err != nil {
		return err
	}
	return nil
}
func (e *Engine) rawKey(c herbcache.Context, key []byte) []byte {
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
func (e *Engine) loadVersion(key []byte, cacheable bool) ([]byte, error) {
	if cacheable {
		return e.getCachedVersion(key)
	}
	return e.LoadRawVersion(key)
}
func (e *Engine) getCachedVersion(key []byte) ([]byte, error) {
	version, err := e.Cache.Get(key)
	if err == nil {
		return version, nil
	}
	if err != herbdata.ErrNotFound {
		return nil, err
	}
	version, err = e.LoadRawVersion(key)
	if err != nil {
		return nil, err
	}
	err = e.Cache.SetWithTTL(key, version, e.VersionTTL)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (e *Engine) LoadRawVersion(key []byte) ([]byte, error) {
	v, err := e.VersionStore.Get(key)
	if err == nil {
		return v, nil
	}
	if err == herbdata.ErrNotFound {
		return []byte{}, nil
	}
	return nil, err
}
func (e *Engine) getRawkeyAndVersion(c herbcache.Context, key []byte) ([]byte, []byte, error) {
	var err error
	versionbuf := bytes.NewBuffer(nil)
	keybuf := bytes.NewBuffer(nil)
	_, err = herbcache.WriteNamespace(keybuf, c.Namespace())
	if err != nil {
		panic(err)
	}
	cacheable := e.VersionTTL > 0 && e.VersionStore != nil
	d := c.Position().RootDirectory()
	for d != nil {
		_, err = herbcache.WriteDirectory(keybuf, d)
		if err != nil {
			return nil, nil, err
		}
		currentkey := keybuf.Bytes()
		v, err := e.loadVersion(currentkey, cacheable)
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
func (e *Engine) setVersion(c herbcache.Context, version []byte) error {
	var err error
	cacheable := e.VersionTTL > 0 && e.VersionStore != nil
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
	err = e.VersionStore.Set(key, version)
	if err != nil {
		return err
	}
	if cacheable {
		return e.Cache.Delete(key)
	}
	return nil
}
func (e *Engine) ExecuteGet(c herbcache.Context, key []byte) ([]byte, error) {
	var data []byte
	var version []byte
	var err error
	var ent *enity
	var rawkey []byte
	flushable := c.Flushable()
	if flushable {
		rawkey, version, err = e.getRawkeyAndVersion(c, key)
		if err != nil {
			return nil, err
		}
	} else {
		rawkey = e.rawKey(c, key)
	}
	data, err = e.Cache.Get(rawkey)
	if err != nil {
		return nil, err
	}
	ent, err = loadEnity(data, flushable, version)
	if err != nil {
		if err == ErrEnityTypecodeNotMatch || err == ErrEnityVersionNotMatch {
			return nil, herbdata.ErrNotFound
		}
		return nil, err
	}
	return ent.data, nil
}
func (e *Engine) ExecuteSetWithTTL(c herbcache.Context, key []byte, data []byte, ttl int64) error {
	var version []byte
	var err error
	var ent *enity
	var rawkey []byte
	flushable := c.Flushable()
	if flushable {
		rawkey, version, err = e.getRawkeyAndVersion(c, key)
		if err != nil {
			return err
		}
	} else {
		rawkey = e.rawKey(c, key)
	}
	ent = createEnity(flushable, version, data)
	buf := bytes.NewBuffer(nil)
	err = ent.SaveTo(buf)
	if err != nil {
		return err
	}
	return e.Cache.SetWithTTL(rawkey, buf.Bytes(), ttl)
}
func (e *Engine) ExecuteDelete(c herbcache.Context, key []byte) error {
	return e.Cache.Delete(e.rawKey(c, key))
}
func (e *Engine) ExecuteFlush(c herbcache.Context) error {
	if !c.Flushable() {
		return herbdata.ErrNotFlushable
	}
	if e.VersionStore == nil {
		return ErrNoVersionStore
	}
	v, err := e.VersionGenerator()
	if err != nil {
		return err
	}
	return e.setVersion(c, []byte(v))

}

func New() *Engine {
	return &Engine{
		VersionGenerator: DefaultVersionGenerator,
	}
}
