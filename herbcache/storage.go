package herbcache

type Storage interface {
	ExecuteGet(p Parameter, key []byte) ([]byte, error)
	ExecuteSetWithTTL(p Parameter, key []byte, data []byte, ttl int64) error
	ExecuteDelete(p Parameter, key []byte) error
	ExecuteFlush(p Parameter) error
}
