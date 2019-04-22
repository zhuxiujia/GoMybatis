package GoMybatis

import (
	"sync"
)

type GoroutineSessionMap struct {
	m    map[int64]Session
	lock sync.RWMutex
}

func (it GoroutineSessionMap) New() GoroutineSessionMap {
	return GoroutineSessionMap{
		m: make(map[int64]Session),
	}
}
func (it *GoroutineSessionMap) Put(k int64, session Session) {
	it.lock.Lock()
	defer it.lock.Unlock()
	it.m[k] = session
}
func (it *GoroutineSessionMap) Get(k int64) Session {
	return it.m[k]
}
