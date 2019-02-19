package GoExpress

import (
	"fmt"
	"testing"
)

func TestNode_Run(t *testing.T) {
	var express = "a == 1 && a != 0"

	//express = "1 + 2 > 3 + 6"
	//express = "1 + 2 != nil"

	var node, e = Parser(express)
	if e != nil {
		t.Fatal(e)
	}
	v, e := node.Eval(map[string]interface{}{"a": 1})
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(v)
}

func BenchmarkArgNode_Eval(b *testing.B) {
	b.StopTimer()
	var express = "a != nil"
	var node, e = Parser(express)
	if e != nil {
		b.Fatal(e)
	}
	var m = map[string]interface{}{"a": 2}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, e := node.Eval(m)
		if e != nil {
			b.Fatal(e)
		}
	}
}
