package GoMybatis

import (
	"sync"
)

type GoroutineSessionMap struct {
	m sync.Map
}

func (it GoroutineSessionMap) New() GoroutineSessionMap {
	return GoroutineSessionMap{
		m: sync.Map{},
	}
}
func (it *GoroutineSessionMap) Put(k int64, session Session) {
	it.m.Store(k, session)
}
func (it *GoroutineSessionMap) Get(k int64) Session {
	var v, ok = it.m.Load(k)
	if ok {
		return v.(Session)
	} else {
		return nil
	}
}

func (it *GoroutineSessionMap) Delete(k int64) {
	it.m.Delete(k)
}
