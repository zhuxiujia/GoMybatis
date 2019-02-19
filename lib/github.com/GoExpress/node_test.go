package GoExpress

import (
	"fmt"
	"testing"
)

func TestNode_Run(t *testing.T) {
	var express = "1 + 2 + 3 + 2 + 3 + 2 + 3 + 2 + 3 + 2 + 3"

	//express = "1 + 2 > 3 + 6"
	//express = "1 + 2 != nil"

	var node, e = Parser(express)
	if e != nil {
		t.Fatal(e)
	}
	v, e := node.Eval(nil)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(v)
}

func BenchmarkArgNode_Eval(b *testing.B) {
	b.StopTimer()
	var express = "1 == 1 && 1 == 1 && 1 == 1 && 1 == 1 && 1 == 1"
	var node, e = Parser(express)
	if e != nil {
		b.Fatal(e)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, e := node.Eval(nil)
		if e != nil {
			b.Fatal(e)
		}
	}
}
