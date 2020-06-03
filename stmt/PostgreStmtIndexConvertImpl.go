package stmt

import (
	"fmt"
	"sync"
)

type PostgreStmtIndexConvertImpl struct {
	sync.RWMutex
	Counter int
}

func (p *PostgreStmtIndexConvertImpl) Inc() {
	p.Lock()
	defer p.Unlock()
	p.Counter++
}

func (p *PostgreStmtIndexConvertImpl) Get() int {
	p.RLock()
	defer p.RUnlock()
	return p.Counter
}

func (p *PostgreStmtIndexConvertImpl) Convert() string {
	return fmt.Sprint(" $", p.Get(), " ")
}
