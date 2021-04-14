package cachecommand

import (
	"github.com/herb-go/herbdata/dataencoding"

	"github.com/herb-go/datamodules/herbcache"
)

type OperationCode byte

func (c OperationCode) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideOperationCode(c), nil
}

const OperationCodeSetWithTTL = OperationCode(1)
const OperationCodeGet = OperationCode(2)
const OperationCodeDelete = OperationCode(3)
const OperationCodeFlush = OperationCode(4)

type Context struct {
	preset        *Preset
	ttl           int64
	cache         *herbcache.Cache
	encoding      *dataencoding.Encoding
	operationCode OperationCode
	key           []byte
	data          []byte
	loader        func([]byte) ([]byte, error)
}

func (c *Context) Loader() func([]byte) ([]byte, error) {
	return c.loader
}

func (c *Context) OverrideLoader(loader func([]byte) ([]byte, error)) *Context {
	ctx := c.Clone()
	ctx.loader = loader
	return ctx
}

func (c *Context) OperationCode() OperationCode {
	return c.operationCode
}
func (c *Context) OverrideOperationCode(code OperationCode) *Context {
	ctx := c.Clone()
	ctx.operationCode = code
	return ctx
}
func (c *Context) Encoding() *dataencoding.Encoding {
	return c.encoding
}
func (c *Context) OverrideEncoding(encoding *dataencoding.Encoding) *Context {
	ctx := c.Clone()
	ctx.encoding = encoding
	return ctx
}
func (c *Context) Cache() *herbcache.Cache {
	return c.cache
}
func (c *Context) OverrideCache(cache *herbcache.Cache) *Context {
	ctx := c.Clone()
	ctx.cache = cache
	return ctx
}
func (c *Context) TTL() int64 {
	return c.ttl
}
func (c *Context) OverrideTTL(ttl int64) *Context {
	ctx := c.Clone()
	ctx.ttl = ttl
	return ctx
}
func (c *Context) Data() []byte {
	return c.data
}
func (c *Context) OverrideData(data []byte) *Context {
	ctx := c.Clone()
	ctx.data = data
	return ctx
}
func (c *Context) Key() []byte {
	return c.key
}
func (c *Context) OverrideKey(key []byte) *Context {
	ctx := c.Clone()
	ctx.key = key
	return ctx
}
func (c *Context) Clone() *Context {
	return &Context{
		preset:        c.preset,
		ttl:           c.ttl,
		cache:         c.cache,
		encoding:      c.encoding,
		operationCode: c.operationCode,
		key:           c.key,
		data:          c.data,
	}
}

func (c *Context) Concat(cmds ...Command) *Context {
	ctx := c.Clone()
	ctx.preset = ctx.preset.Concat(cmds...)
	return ctx
}

func (c *Context) OverrideCacheFlushable(flashable bool) *Context {
	return c.Concat(Flushable(flashable))
}

func (c *Context) Allocate(name string) *Context {
	return c.Concat(Allocate(name))
}

func (c *Context) ChildCache(name string) *Context {
	return c.Concat(ChildCache(name))
}

func (c *Context) PrefixCache(prefix string) *Context {
	return c.Concat(PrefixCache(prefix))
}
func (c *Context) Flush() error {
	_, err := c.Concat(OperationCodeFlush, Operate).Exec()
	return err
}
func (c *Context) Delete(key []byte) error {
	_, err := c.Concat(Key(key), OperationCodeDelete, Operate).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (c *Context) SDelete(key string) error {
	return c.Delete([]byte(key))
}
func (c *Context) Get(key []byte) ([]byte, error) {
	ctx, err := c.Concat(Key(key), OperationCodeGet, Operate).Exec()
	if err != nil {
		return nil, err
	}
	return ctx.data, nil
}
func (c *Context) SGet(key []byte) ([]byte, error) {
	return c.Get([]byte(key))
}
func (c *Context) SetWithTTL(key []byte, data []byte, ttl int64) error {
	_, err := c.Concat(Key(key), Data(data), TTL(ttl), OperationCodeSetWithTTL, Operate).Exec()
	return err
}
func (c *Context) SSetWithTTL(key string, data []byte, ttl int64) error {
	return c.SetWithTTL([]byte(key), data, ttl)
}
func (c *Context) Exec() (*Context, error) {
	return c.preset.Exec(c)
}

func NewContext() *Context {
	return &Context{}
}
