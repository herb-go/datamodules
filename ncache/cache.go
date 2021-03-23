package ncache

import (
	"bytes"

	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/datautil"
)

func Join(pathlist ...[]byte) []byte {
	buf := bytes.NewBuffer(nil)
	err := datautil.PackTo(buf, nil, pathlist...)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type CachePathPrefix []byte

func (p CachePathPrefix) MustJoin(pathlist ...[]byte) []byte {
	buf := bytes.NewBuffer(nil)
	_, err := buf.Write(p)
	if err != nil {
		panic(err)
	}
	err = datautil.PackTo(buf, nil, pathlist...)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

var CachePathPrefixValue = CachePathPrefix([]byte{0})
var CachePathPrefixVersion = CachePathPrefix([]byte{1})

type Cache struct {
	revocable     bool
	namespaceTree [][]byte
	storage       *Storage
	todos         []Directive
}

func (c *Cache) Revocable() bool {
	return c.revocable
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
func (c *Cache) getVersion() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var err error
	cacheable := c.storage.VersionTTL > 0 && c.storage.VersionStore != nil
	for k := range c.namespaceTree {
		var version []byte
		key := CachePathPrefixVersion.MustJoin(c.namespaceTree[0 : k+1]...)
		if cacheable {
			version, err = c.getCachedVersion(key)
		} else {
			version, err = c.storage.LoadRawVersion(key)
		}
		if err != nil {
			return nil, err
		}
		err = datautil.PackTo(buf, nil, version)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
func (c *Cache) setVersion(version []byte) error {
	cacheable := c.storage.VersionTTL > 0 && c.storage.VersionStore != nil
	key := CachePathPrefixVersion.MustJoin(c.namespaceTree...)
	err := c.storage.VersionStore.Set(key, version)
	if err != nil {
		return err
	}
	if cacheable {
		return c.storage.Cache.Delete(key)
	}
	return nil
}
func (c *Cache) mustGetNamespace() []byte {
	return Join(c.namespaceTree...)
}
func (c *Cache) Revoke() error {
	if !c.revocable {
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
	namespace := c.mustGetNamespace()
	if c.revocable {
		version, err = c.getVersion()
		if err != nil {
			return nil, err
		}
	}
	data, err = c.storage.Cache.Get(CachePathPrefixValue.MustJoin(namespace, key))
	if err != nil {
		return nil, err
	}
	e, err = loadEnity(data, c.revocable, version)
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
	namespace := c.mustGetNamespace()
	if c.revocable {
		version, err = c.getVersion()
		if err != nil {
			return err
		}
	}
	e = createEnity(c.revocable, version, data)
	buf := bytes.NewBuffer(nil)
	err = e.SaveTo(buf)
	if err != nil {
		return err
	}
	return c.storage.Cache.SetWithTTL(CachePathPrefixValue.MustJoin(namespace, key), buf.Bytes(), ttl)
}

func (c *Cache) Delete(key []byte) error {
	namespace := c.mustGetNamespace()
	return c.storage.Cache.Delete(CachePathPrefixValue.MustJoin(namespace, key))
}

func (c *Cache) Clone() *Cache {
	t := make([][]byte, len(c.namespaceTree))
	for k := range t {
		t[k] = append([]byte{}, c.namespaceTree[k]...)
	}
	return &Cache{
		revocable:     c.revocable,
		storage:       c.storage,
		namespaceTree: t,
		todos:         append([]Directive{}, c.todos...),
	}
}
func (c *Cache) NamescapedCache(namescape []byte) herbdata.NestableCache {
	return c.VaryNamesapce(namescape)
}
func (c *Cache) ChildCache(name []byte) herbdata.NestableCache {
	return c.Child(name)
}

func (c *Cache) VaryRevocable(revocable bool) *Cache {
	cc := c.Clone()
	cc.revocable = revocable
	return cc
}
func (c *Cache) buildNamespace(prefix []byte, suffixs ...[]byte) {
	var err error
	buf := bytes.NewBuffer(nil)
	if len(prefix) > 0 {
		_, err = buf.Write(prefix)
		if err != nil {
			panic(err)
		}
	}
	err = datautil.PackTo(buf, nil, suffixs...)
	if err != nil {
		panic(err)
	}
	index := len(c.namespaceTree) - 1
	c.namespaceTree[index] = buf.Bytes()
}
func (c *Cache) VarySuffix(suffixs ...[]byte) *Cache {
	index := len(c.namespaceTree) - 1
	cc := c.Clone()
	cc.buildNamespace(c.namespaceTree[index], suffixs...)
	return cc
}
func (c *Cache) VaryNamesapce(namespace ...[]byte) *Cache {
	cc := c.Clone()
	cc.buildNamespace(nil, namespace...)
	return cc
}
func (c *Cache) VaryTodos(todos ...Directive) *Cache {
	cc := c.Clone()
	c.todos = append(c.todos, todos...)
	return cc
}
func (c *Cache) Todos() []Directive {
	return c.todos
}

func (c *Cache) ExecuteTodos() error {
	for len(c.todos) > 0 {
		err := c.todos[0].Execute(c)
		if err != nil {
			return err
		}
		c.todos = c.todos[1:]
	}
	return nil
}

func (c *Cache) Child(name ...[]byte) *Cache {
	cc := c.Clone()
	cc.namespaceTree = append(cc.namespaceTree, name...)
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
		return dst == nil && c == nil
	}
	if len(c.namespaceTree) != len(dst.namespaceTree) {
		return false
	}
	for k := range c.namespaceTree {
		if bytes.Compare(c.namespaceTree[k], dst.namespaceTree[k]) != 0 {
			return false
		}
	}
	return c.revocable == dst.revocable &&
		c.storage == dst.storage
}

func New() *Cache {
	return &Cache{
		namespaceTree: [][]byte{[]byte{}},
	}
}
