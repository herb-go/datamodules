package ncache

import (
	"bytes"
	"io"

	"github.com/herb-go/herbdata/datautil"

	"github.com/herb-go/herbdata"
)

type Cache struct {
	flushable bool
	namespace []byte
	group     []byte
	path      *Path
	storage   *Storage
	promises  []Directive
}

func (c *Cache) Flushable() bool {
	return c.flushable
}
func (c *Cache) writeNamespace(w io.Writer) error {
	var err error
	_, err = writeTokenAndData(w, tokenBeforeNamespace, c.namespace)
	return err
}
func (c *Cache) writePath(w io.Writer) error {
	var err error
	_, err = c.path.WriteTo(w)
	return err
}
func (c *Cache) writeKey(w io.Writer, key []byte) error {
	var err error
	_, err = writeTokenAndData(w, tokenBeforeGroup, c.group)
	if err != nil {
		return err
	}
	_, err = writeTokenAndData(w, tokenBeforeValue, key)
	return err
}
func (c *Cache) rawKey(key []byte) []byte {
	var err error
	buf := bytes.NewBuffer(nil)
	err = c.writeNamespace(buf)
	if err != nil {
		panic(err)
	}
	err = c.writePath(buf)
	if err != nil {
		panic(err)
	}
	err = c.writeKey(buf, key)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (c *Cache) loadVersion(key []byte, cacheable bool) ([]byte, error) {
	if cacheable {
		return c.getCachedVersion(key)
	}
	return c.storage.LoadRawVersion(key)
}
func (c *Cache) getRawkeyAndVersion(key []byte) ([]byte, []byte, error) {
	var err error
	versionbuf := bytes.NewBuffer(nil)
	keybuf := bytes.NewBuffer(nil)
	err = c.writeNamespace(keybuf)
	if err != nil {
		panic(err)
	}
	cacheable := c.storage.VersionTTL > 0 && c.storage.VersionStore != nil
	np := c.path.toNextPath(nil)
	for np != nil {
		_, err = np.writeSelf(keybuf)
		if err != nil {
			return nil, nil, err
		}
		currentkey := keybuf.Bytes()
		v, err := c.loadVersion(currentkey, cacheable)
		if err != nil {
			return nil, nil, err
		}
		err = datautil.PackTo(versionbuf, nil, v)
		if err != nil {
			return nil, nil, err
		}
		keybuf = bytes.NewBuffer(currentkey)
		np = np.next
	}
	err = c.writeKey(keybuf, key)
	if err != nil {
		return nil, nil, err
	}
	return keybuf.Bytes(), versionbuf.Bytes(), nil
}
func (c *Cache) getCachedVersion(key []byte) ([]byte, error) {
	version, err := c.storage.Cache.Get(key)
	if err == nil {
		return version, nil
	}
	if err != herbdata.ErrNotFound {
		return nil, err
	}
	version, err = c.storage.LoadRawVersion(key)
	if err != nil {
		return nil, err
	}
	err = c.storage.Cache.SetWithTTL(key, version, c.storage.VersionTTL)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (c *Cache) setVersion(version []byte) error {
	var err error
	cacheable := c.storage.VersionTTL > 0 && c.storage.VersionStore != nil
	keybuf := bytes.NewBuffer(nil)
	err = c.writeNamespace(keybuf)
	if err != nil {
		return err
	}
	_, err = c.path.toNextPath(nil).writeAll(keybuf)
	if err != nil {
		return err
	}
	key := keybuf.Bytes()
	err = c.storage.VersionStore.Set(key, version)
	if err != nil {
		return err
	}
	if cacheable {
		return c.storage.Cache.Delete(key)
	}
	return nil
}

func (c *Cache) Flush() error {
	if !c.flushable {
		return herbdata.ErrIrrevocable
	}
	if c.storage.VersionStore == nil {
		return ErrNoVersionStore
	}
	v, err := c.storage.VersionGenerator()
	if err != nil {
		return err
	}
	return c.setVersion([]byte(v))
}
func (c *Cache) Get(key []byte) ([]byte, error) {
	var data []byte
	var version []byte
	var err error
	var e *enity
	var rawkey []byte
	if c.flushable {
		rawkey, version, err = c.getRawkeyAndVersion(key)
		if err != nil {
			return nil, err
		}
	} else {
		rawkey = c.rawKey(key)
	}
	data, err = c.storage.Cache.Get(rawkey)
	if err != nil {
		return nil, err
	}
	e, err = loadEnity(data, c.flushable, version)
	if err != nil {
		if err == ErrEnityTypecodeNotMatch || err == ErrEnityVersionNotMatch {
			return nil, herbdata.ErrNotFound
		}
		return nil, err
	}
	return e.data, nil

}

func (c *Cache) SetWithTTL(key []byte, data []byte, ttl int64) error {
	var version []byte
	var err error
	var e *enity
	var rawkey []byte
	if c.flushable {
		rawkey, version, err = c.getRawkeyAndVersion(key)
		if err != nil {
			return err
		}
	} else {
		rawkey = c.rawKey(key)
	}
	e = createEnity(c.flushable, version, data)
	buf := bytes.NewBuffer(nil)
	err = e.SaveTo(buf)
	if err != nil {
		return err
	}
	return c.storage.Cache.SetWithTTL(rawkey, buf.Bytes(), ttl)
}

func (c *Cache) Delete(key []byte) error {
	return c.storage.Cache.Delete(c.rawKey(key))
}

func (c *Cache) Clone() *Cache {
	return &Cache{
		flushable: c.flushable,
		storage:   c.storage,
		path:      c.path,
		group:     c.group,
		namespace: c.namespace,
		promises:  append([]Directive{}, c.promises...),
	}
}

func (c *Cache) SubCache(name []byte) herbdata.NestableCache {
	return c.Child(name)
}

func (c *Cache) VaryFlushable(flushable bool) *Cache {
	cc := c.Clone()
	cc.flushable = flushable
	return cc
}

func (c *Cache) VaryPrefix(group []byte) *Cache {
	cc := c.Clone()
	SetCacheGroup(cc, group)
	return cc
}
func (c *Cache) Migrate(namespace []byte) *Cache {
	cc := c.Clone()
	SetCacheNamespace(cc, namespace)
	SetCachePath(cc, nil)
	SetCacheGroup(cc, nil)
	return cc
}

func (c *Cache) VaryMorePromises(promises ...Directive) *Cache {
	cc := c.Clone()
	c.promises = append(c.promises, promises...)
	return cc
}
func (c *Cache) Promises() []Directive {
	return c.promises
}

func (c *Cache) ResolvePromises() error {
	p := c.promises
	c.promises = nil
	for len(p) > 0 {
		err := c.promises[0].Execute(c)
		if err != nil {
			return err
		}
		c.promises = c.promises[1:]
	}
	return nil
}

func (c *Cache) Child(name []byte) *Cache {
	cc := c.Clone()
	SetCachePath(cc, cc.path.Append(c.group, name))
	SetCacheGroup(cc, nil)
	return cc
}
func (c *Cache) VaryStorage(storage *Storage) *Cache {
	cc := c.Clone()
	cc.storage = storage
	return cc
}
func (c *Cache) Storage() *Storage {
	return c.storage
}

func (c *Cache) CopyFrom(src *Cache) {
	SetCache(c, src.Clone())
}
func (c *Cache) Equal(dst *Cache) bool {
	if dst == nil || c == nil {
		return dst == c
	}
	if !c.path.Equal(dst.path) {
		return false
	}
	if !bytes.Equal(c.group, dst.group) {
		return false
	}
	if !bytes.Equal(c.namespace, dst.namespace) {
		return false
	}
	return c.flushable == dst.flushable &&
		c.storage == dst.storage
}

func New() *Cache {
	return &Cache{}
}
