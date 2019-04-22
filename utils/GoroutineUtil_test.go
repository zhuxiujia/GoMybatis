package utils

import "testing"

func TestGoroutineID(t *testing.T) {
	println(GoroutineID())
}

func BenchmarkGoroutineID(b *testing.B) {
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GoroutineID()
	}
	//BenchmarkGoroutineID-8   	1000000000	         2.91 ns/op
}
