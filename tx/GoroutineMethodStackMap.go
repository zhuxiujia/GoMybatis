package tx

import (
	"sync"
)

type GoroutineMethodStackMap struct {
	m    map[int64]*StructField
	lock sync.RWMutex
}

func (it GoroutineMethodStackMap) New() GoroutineMethodStackMap {
	return GoroutineMethodStackMap{
		m: make(map[int64]*StructField),
	}
}
func (it *GoroutineMethodStackMap) Put(k int64, methodInfo *StructField) {
	it.lock.Lock()
	defer it.lock.Unlock()
	it.m[k] = methodInfo
}
func (it *GoroutineMethodStackMap) Get(k int64) *StructField {
	return it.m[k]
}
