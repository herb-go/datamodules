package cachepreset

import (
	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/dataencoding"

	"github.com/herb-go/datamodules/herbcache"
)

type Preset struct {
	commands *Commands
	ttl      int64
	cache    *herbcache.Cache
	encoding *dataencoding.Encoding
	key      []byte
	lockers  *Lockers
	data     []byte
	loader   func([]byte) ([]byte, error)
}

func (p *Preset) Exec(preset *Preset) (*Preset, error) {
	newpreset := p.Clone()
	newpreset.commands = preset.commands.Concat(p.commands)
	return newpreset, nil
}

func (p *Preset) Lockers() *Lockers {
	return p.lockers
}

func (p *Preset) OverrideLockers(lockers *Lockers) *Preset {
	preset := p.Clone()
	preset.lockers = lockers
	return preset
}
func (p *Preset) Loader() func([]byte) ([]byte, error) {
	return p.loader
}

func (p *Preset) OverrideLoader(loader func([]byte) ([]byte, error)) *Preset {
	preset := p.Clone()
	preset.loader = loader
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
		commands: p.commands,
		ttl:      p.ttl,
		cache:    p.cache,
		encoding: p.encoding,
		key:      p.key,
		data:     p.data,
		loader:   p.loader,
	}
}

func (p *Preset) Concat(cmds ...Command) *Preset {
	preset := p.Clone()
	preset.commands = preset.commands.Concat(cmds...)
	return preset
}

func (p *Preset) Flushable(flashable bool) *Preset {
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
	preset, err := p.Apply()
	if err != nil {
		return err
	}
	return preset.cache.Flush()
}
func (p *Preset) Delete(key []byte) error {
	preset, err := p.Concat(Key(key)).Apply()
	if err != nil {
		return err
	}
	return p.cache.Delete(preset.key)
}
func (p *Preset) DeleteS(key string) error {
	return p.Delete([]byte(key))
}
func (p *Preset) Get(key []byte) ([]byte, error) {
	preset, err := p.get(key)
	if err != nil {
		return nil, err
	}
	return preset.data, nil
}
func (p *Preset) get(key []byte) (*Preset, error) {
	preset, err := p.Concat(Key(key)).Apply()
	if err != nil {
		return nil, err
	}
	data, err := preset.cache.Get(preset.key)
	if err == nil {
		preset.data = data
		return preset, nil
	}
	if err != herbdata.ErrNotFound {
		return nil, err
	}
	if preset.loader == nil {
		return nil, err
	}
	var blocked bool
	var needunlock bool
	if preset.lockers != nil {
		blocked = !preset.lockers.m.Lock(string(key))
		needunlock = true
	}
	defer func() {
		if needunlock {
			preset.lockers.m.Unlock(string(key))
		}
	}()
	//retry read cache if blocker
	if blocked {
		data, err = preset.cache.Get(preset.key)
		if err == nil {
			preset.data = data
			return preset, nil
		}
		if err != herbdata.ErrNotFound {
			return nil, err
		}
	}
	data, err = preset.loader(preset.key)
	if err != nil {
		return nil, err
	}
	if preset.lockers != nil {
		preset.lockers.m.Unlock(string(key))
		needunlock = false
	}
	if preset.ttl > 0 {
		err = preset.cache.SetWithTTL(preset.key, preset.data, preset.ttl)
		if err != nil {
			return nil, err
		}
	}
	preset.data = data
	return preset, nil
}
func (p *Preset) GetS(key []byte) ([]byte, error) {
	return p.Get([]byte(key))
}
func (p *Preset) SetWithTTL(key []byte, data []byte, ttl int64) error {
	preset, err := p.Concat(Key(key), Data(data), TTL(ttl)).Apply()
	if err != nil {
		return err
	}
	return preset.cache.SetWithTTL(preset.key, preset.data, preset.ttl)
}
func (p *Preset) SetWithTTLS(key string, data []byte, ttl int64) error {
	return p.SetWithTTL([]byte(key), data, ttl)
}
func (p *Preset) Load(key []byte, v interface{}) error {
	preset, err := p.get(key)
	if err != nil {
		return err
	}
	return preset.encoding.Unmarshal(preset.data, v)
}
func (p *Preset) LoadS(key string, v interface{}) error {
	return p.Load([]byte(key), v)
}
func (p *Preset) Apply() (*Preset, error) {
	return p.commands.Exec(p)
}

func New(commands ...Command) *Preset {
	p := &Preset{
		encoding: dataencoding.NopEncoding,
	}
	return p.Concat(commands...)
}
