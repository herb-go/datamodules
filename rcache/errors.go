package rcache

import "errors"

var ErrUnresolvableCacheEnity = errors.New("unresolvable cache enity")
var ErrEnityTypecodeNotMatch = errors.New("enity typecode not match")
var ErrEnityVersionNotMatch = errors.New("enity version not match")
var ErrCacheIrrevocable = errors.New("cache irrevocable")
