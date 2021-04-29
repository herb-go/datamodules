package cachepreset

import (
	"github.com/herb-go/herbdata/dataencoding"

	"github.com/herb-go/datamodules/herbcache"
)

type OperationCode byte

func (c OperationCode) Exec(preset *Preset) (newpreset *Preset, err error) {
	return preset.OverrideOperationCode(c), nil
}

const OperationCodeSetWithTTL = OperationCode(1)
const OperationCodeGet = OperationCode(2)
const OperationCodeDelete = OperationCode(3)
const OperationCodeFlush = OperationCode(4)

type Preset struct {
	commands      *Commands
	ttl           int64
	cache         *herbcache.Cache
	encoding      *dataencoding.Encoding
	operationCode OperationCode
	key           []byte
	data          []byte
	loader        func([]byte) ([]byte, error)
}

func (p *Preset) Loader() func([]byte) ([]byte, error) {
	return p.loader
}

func (p *Preset) OverrideLoader(loader func([]byte) ([]byte, error)) *Preset {
	preset := p.Clone()
	preset.loader = loader
	return preset
}

func (p *Preset) OperationCode() OperationCode {
	return p.operationCode
}
func (p *Preset) OverrideOperationCode(code OperationCode) *Preset {
	preset := p.Clone()
	preset.operationCode = code
	return preset
}
func (p *Preset) Encoding() *dataencoding.Encoding {
	return p.encoding
}
func (p *Preset) OverrideEncoding(encoding *dataencoding.Encoding) *Preset {
	preset := p.Clone()
	preset.encoding = encoding
	return preset
}
func (p *Preset) Cache() *herbcache.Cache {
	return p.cache
}
func (p *Preset) OverrideCache(cache *herbcache.Cache) *Preset {
	preset := p.Clone()
	preset.cache = cache
	return preset
}
func (p *Preset) TTL() int64 {
	return p.ttl
}
func (p *Preset) OverrideTTL(ttl int64) *Preset {
	preset := p.Clone()
	preset.ttl = ttl
	return preset
}
func (p *Preset) Data() []byte {
	return p.data
}
func (p *Preset) OverrideData(data []byte) *Preset {
	preset := p.Clone()
	preset.data = data
	return preset
}
func (p *Preset) Key() []byte {
	return p.key
}
func (p *Preset) OverrideKey(key []byte) *Preset {
	preset := p.Clone()
	preset.key = key
	return preset
}
func (p *Preset) Clone() *Preset {
	return &Preset{
		commands:      p.commands,
		ttl:           p.ttl,
		cache:         p.cache,
		encoding:      p.encoding,
		operationCode: p.operationCode,
		key:           p.key,
		data:          p.data,
	}
}

func (p *Preset) Concat(cmds ...Command) *Preset {
	preset := p.Clone()
	preset.commands = preset.commands.Concat(cmds...)
	return preset
}

func (p *Preset) OverrideCacheFlushable(flashable bool) *Preset {
	return p.Concat(Flushable(flashable))
}

func (p *Preset) Allocate(name string) *Preset {
	return p.Concat(Allocate(name))
}

func (p *Preset) ChildCache(name string) *Preset {
	return p.Concat(ChildCache(name))
}

func (p *Preset) PrefixCache(prefix string) *Preset {
	return p.Concat(PrefixCache(prefix))
}
func (p *Preset) Flush() error {
	_, err := p.Concat(OperationCodeFlush, Operate).Exec()
	return err
}
func (p *Preset) Delete(key []byte) error {
	_, err := p.Concat(Key(key), OperationCodeDelete, Operate).Exec()
	if err != nil {
		return err
	}
	return nil
}
func (p *Preset) SDelete(key string) error {
	return p.Delete([]byte(key))
}
func (p *Preset) Get(key []byte) ([]byte, error) {
	preset, err := p.Concat(Key(key), OperationCodeGet, Operate).Exec()
	if err != nil {
		return nil, err
	}
	return preset.data, nil
}
func (p *Preset) SGet(key []byte) ([]byte, error) {
	return p.Get([]byte(key))
}
func (p *Preset) SetWithTTL(key []byte, data []byte, ttl int64) error {
	_, err := p.Concat(Key(key), Data(data), TTL(ttl), OperationCodeSetWithTTL, Operate).Exec()
	return err
}
func (p *Preset) SSetWithTTL(key string, data []byte, ttl int64) error {
	return p.SetWithTTL([]byte(key), data, ttl)
}
func (p *Preset) Exec() (*Preset, error) {
	return p.commands.Exec(p)
}

func NewContext() *Preset {
	return &Preset{}
}
