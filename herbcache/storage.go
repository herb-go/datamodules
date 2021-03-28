package herbcache

type Storage interface {
	ExecuteGet(c Context, key []byte) ([]byte, error)
	ExecuteSetWithTTL(c Context, key []byte, data []byte, ttl int64) error
	ExecuteDelete(c Context, key []byte) error
	ExecuteFlush(c Context) error
}

type NopStorage struct{}

func (s *NopStorage) ExecuteGet(c Context, key []byte) ([]byte, error) {
	return nil, ErrStorageRequired
}
func (s *NopStorage) ExecuteSetWithTTL(c Context, key []byte, data []byte, ttl int64) error {
	return ErrStorageRequired
}
func (s *NopStorage) ExecuteDelete(c Context, key []byte) error {
	return ErrStorageRequired
}
func (s *NopStorage) ExecuteFlush(c Context) error {
	return ErrStorageRequired
}
