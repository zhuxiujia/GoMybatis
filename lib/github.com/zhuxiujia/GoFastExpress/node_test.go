package GoFastExpress

import (
	"fmt"
	"testing"
)

func TestNode_Run(t *testing.T) {

	var expressions = []string{
		"`f`+`s`",
		"a +1 > b * 8",
		"a >= 0",
		"'a'+c",
		"b",
		"a < 1",
		"a +1 > b*8",
		"a * b == 2",
		"a - b == 0",
		"a >= 0 && a != 0",
		"a == 1 && a != 0",
		"1 > 3 ",
		"1 + 2 != nil",
		"1 != null",
		"1 + 2 != nil && 1 > 0 ",
		"1 + 2 != nil && 2 < b*8 ",
	}
	for _, expr := range expressions {
		node, e := Parser(expr)
		if e != nil {
			t.Fatal(e)
		}
		v, e := node.Eval(map[string]interface{}{"a": 1, "b": 2, "c": "c"})
		if e != nil {
			t.Fatal(e)
		}
		fmt.Println(v)
	}
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

type As struct {
	B string
}

func TestArgNode_Eval_Take(b *testing.T) {
	var express = "a.B"
	var node, e = Parser(express)
	if e != nil {
		b.Fatal(e)
	}
	var m = map[string]interface{}{"a": As{B: "sdffdasf"}}
	r, e := node.Eval(m)
	if e != nil {
		b.Fatal(e)
	}
	fmt.Println(r)
}

func BenchmarkArgNode_Eval_Take(b *testing.B) {
	b.StopTimer()
	var express = "a.B"
	var node, e = Parser(express)
	if e != nil {
		b.Fatal(e)
	}
	var m = map[string]interface{}{"a": As{B: "sdffdasf"}}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, e := node.Eval(m)
		if e != nil {
			b.Fatal(e)
		}
		//fmt.Println(r)
	}
}
