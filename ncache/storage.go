package ncache

import (
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/datautil"
)

type Storage struct {
	VersionGenerator func() (string, error)
	VersionTTL       int64
	VersionStore     herbdata.SetterGetterServer
	Cache            herbdata.CacheServer
}

func (s *Storage) Execute(c *Cache) error {
	c.storage = s
	return nil
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

var DefaultVersionGenerator = func() (string, error) {
	v, err := datautil.Encode(uint64(time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func NewStorage() *Storage {
	return &Storage{
		VersionGenerator: DefaultVersionGenerator,
	}
}
