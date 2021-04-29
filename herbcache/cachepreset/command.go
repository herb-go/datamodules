package cachepreset

import (
	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/herbdata"
)

type Command interface {
	Exec(preset *Preset) (newpreset *Preset, err error)
}

type CommandFunc func(preset *Preset) (newpreset *Preset, err error)

func (f CommandFunc) Exec(preset *Preset) (newpreset *Preset, err error) {
	return f(preset)
}

type Key []byte

func (k Key) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideKey([]byte(k)), nil
}

type Data []byte

func (d Data) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideData([]byte(d)), nil
}

type TTL int64

func (t TTL) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideTTL(int64(t)), nil
}

type Loader func([]byte) ([]byte, error)

func (l Loader) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideLoader(l), nil
}

var Operate = CommandFunc(func(preset *Preset) (newpreset *Preset, err error) {
	c := preset.Clone()
	switch c.operationCode {
	case OperationCodeDelete:
		err = c.cache.Delete(preset.key)
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
		data, err := preset.cache.Get(c.key)
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
	return CommandFunc(func(preset *Preset) (newpreset *Preset, err error) {
		return preset.OverrideCache(cache), nil
	})
}

type Flushable bool

func (f Flushable) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideCache(preset.cache.OverrideFlushable(bool(f))), nil
}

type Allocate string

func (a Allocate) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideCache(preset.cache.Allocate(string(a))), nil
}

type ChildCache string

func (c ChildCache) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideCache(preset.cache.ChildCache(string(c))), nil
}

type PrefixCache string

func (p PrefixCache) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideCache(preset.cache.ChildCache(string(p))), nil
}
