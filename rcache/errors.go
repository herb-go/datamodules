package rcache

import "errors"

var ErrUnresolvedCacheEnity = errors.New("unresolved cache enity")
var ErrEnityTypecodeNotMatch = errors.New("enity typecode not match")
var ErrEnityVersionNotMatch = errors.New("enity version not match")
