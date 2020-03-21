package ids

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetSnowflakeId(t *testing.T) {
	n, _ := NewNode(0)
	id := n.Generate()
	fmt.Println(id)
}

func TestSnowflakeDataRace(t *testing.T) {
	var total = 100000
	n, _ := NewNode(0)
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go func() {
			for i := 0; i < total; i++ {
				n.Generate()
			}
		}()
	}
	for i := 0; i < total; i++ {
		n.Generate()
	}
}

func BenchmarkGetSnowflakeId(b *testing.B) {
	b.StopTimer()
	n, _ := NewNode(0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n.Generate()
	}
}
