package rcache

import "sync"

//Locker cache locker
type Locker struct {
	sync.RWMutex
	Map *sync.Map
	Key []byte
}

//Unlock unlock and delete locker
func (l *Locker) Unlock() {
	l.RWMutex.Unlock()
	l.Map.Delete(l.Key)
}

//Util cache util
type Util struct {
	locks *sync.Map
}

//Locker create new locker with given key.
//Return locker and if locker is locked.
func (u *Util) Locker(key []byte) (*Locker, bool) {
	newlocker := &Locker{
		Map: u.locks,
		Key: key,
	}
	v, ok := u.locks.LoadOrStore(key, newlocker)
	return v.(*Locker), ok
}
