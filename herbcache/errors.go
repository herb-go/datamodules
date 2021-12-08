package herbcache

import "errors"

var ErrStorageRequired = errors.New("herbcache: storage required")
var ErrNotCacheable = errors.New("herbcache: not cacheable")
