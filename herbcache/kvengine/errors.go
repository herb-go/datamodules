package kvengine

import "errors"

var ErrUnresolvableCacheEnity = errors.New("ncache unresolvable cache enity")
var ErrEnityTypecodeNotMatch = errors.New("ncache enity typecode not match")
var ErrEnityVersionNotMatch = errors.New("ncache enity version not match")
var ErrNoVersionStore = errors.New("ncache no version store")
