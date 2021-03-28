package herbcache

type Storage interface {
	ExecuteGet(c Context, key []byte) ([]byte, error)
	ExecuteSetWithTTL(c Context, key []byte, data []byte, ttl int64) error
	ExecuteDelete(c Context, key []byte) error
	ExecuteFlush(c Context) error
}
