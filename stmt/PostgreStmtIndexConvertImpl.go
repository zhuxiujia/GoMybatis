package stmt

import (
	"fmt"
	"sync"
)

type PostgreStmtIndexConvertImpl struct {
	sync.RWMutex
	counter int
}

func (p *PostgreStmtIndexConvertImpl) Inc() {
	p.Lock()
	defer p.Unlock()
	p.counter++
}

func (p *PostgreStmtIndexConvertImpl) Get() int {
	p.RLock()
	defer p.RUnlock()
	return p.counter
}

func (p *PostgreStmtIndexConvertImpl) Convert() string {
	return fmt.Sprint(" $", p.Get(), " ")
}
