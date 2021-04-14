package cachecommand

import (
	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata"
)

type Command interface {
	Exec(ctx *Context) (newctx *Context, err error)
}

type CommandFunc func(ctx *Context) (newctx *Context, err error)

func (f CommandFunc) Exec(ctx *Context) (newctx *Context, err error) {
	return f(ctx)
}

type Key []byte

func (k Key) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideKey([]byte(k)), nil
}

type Data []byte

func (d Data) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideData([]byte(d)), nil
}

type TTL int64

func (t TTL) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideTTL(int64(t)), nil
}

type Loader func([]byte) ([]byte, error)

func (l Loader) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideLoader(l), nil
}

var Operate = CommandFunc(func(ctx *Context) (newctx *Context, err error) {
	c := ctx.Clone()
	switch c.operationCode {
	case OperationCodeDelete:
		err = c.cache.Delete(ctx.key)
		if err != nil {
			return nil, err
		}
	case OperationCodeFlush:
		err = c.cache.Flush()
		if err != nil {
			return nil, err
		}
		return c, nil
	case OperationCodeSetWithTTL:
		err = c.cache.SetWithTTL(c.key, c.data, c.ttl)
		if err != nil {
			return nil, err
		}
		return c, nil
	case OperationCodeGet:
		data, err := ctx.cache.Get(c.key)
		if err == nil {
			c.data = data
			return c, nil
		}
		if c.loader == nil {
			return nil, err
		}
		if err != herbdata.ErrNotFound {
			return nil, err
		}
		data, err = c.loader(c.key)
		if err != nil {
			return nil, err
		}
		if c.ttl != 0 {
			err = c.cache.SetWithTTL(c.key, data, c.ttl)
			if err != nil {
				return nil, err
			}
		}
		return c, nil
	}
	return nil, ErrUnknownOperation
})

func Cache(cache *herbcache.Cache) Command {
	return CommandFunc(func(ctx *Context) (newctx *Context, err error) {
		return ctx.OverrideCache(cache), nil
	})
}

type Flushable bool

func (f Flushable) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideCache(ctx.cache.OverrideFlushable(bool(f))), nil
}

type Allocate string

func (a Allocate) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideCache(ctx.cache.Allocate(string(a))), nil
}

type ChildCache string

func (c ChildCache) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideCache(ctx.cache.ChildCache(string(c))), nil
}

type PrefixCache string

func (p PrefixCache) Exec(ctx *Context) (newctx *Context, err error) {
	return ctx.OverrideCache(ctx.cache.ChildCache(string(p))), nil
}
