package stmt

import (
	"fmt"
	"sync"
)

type OracleStmtIndexConvertImpl struct {
	sync.RWMutex
	counter int
}

func (it *OracleStmtIndexConvertImpl) Convert() string {
	return fmt.Sprint(" :val", it.Get(), " ")
}

func (it *OracleStmtIndexConvertImpl) Inc() {
	it.Lock()
	defer it.Unlock()
	it.counter++
}

func (it *OracleStmtIndexConvertImpl) Get() int {
	it.RLock()
	defer it.RUnlock()
	return it.counter
}
