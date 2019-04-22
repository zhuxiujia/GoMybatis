package goroutineid

import (
	"sync"
	"testing"
)

func TestGoID(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	id := GetGoID()
	go func() {
		idInternal := GetGoID()
		if id == idInternal {
			panic("GoID not success!")
		}
		wg.Done()
	}()
	wg.Wait()
}
