package herbcache

type Engine interface {
	Start() error
	Stop() error
	ExecuteGet(c Context, key []byte) ([]byte, error)
	ExecuteSetWithTTL(c Context, key []byte, data []byte, ttl int64) error
	ExecuteDelete(c Context, key []byte) error
	ExecuteFlush(c Context) error
}
type Storage struct {
	Engine
}

type NopEngine struct{}

func (s *NopEngine) Start() error {
	return ErrStorageRequired
}
func (s *NopEngine) Stop() error {
	return ErrStorageRequired
}
func (s *NopEngine) ExecuteGet(c Context, key []byte) ([]byte, error) {
	return nil, ErrStorageRequired
}
func (s *NopEngine) ExecuteSetWithTTL(c Context, key []byte, data []byte, ttl int64) error {
	return ErrStorageRequired
}
func (s *NopEngine) ExecuteDelete(c Context, key []byte) error {
	return ErrStorageRequired
}
func (s *NopEngine) ExecuteFlush(c Context) error {
	return ErrStorageRequired
}

var DefaultEngine = &NopEngine{}

func NewStorage() *Storage {
	return &Storage{
		Engine: DefaultEngine,
	}
}

type StorageProvider interface {
	Storage() *Storage
}
