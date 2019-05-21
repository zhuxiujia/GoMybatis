package GoFastExpress

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	fmt.Println(Eval("a", "", nil, "a"))
}

func TestEvalNil(t *testing.T) {
	var b *string
	var result, e = Eval("nil!=b", "!=", nil, b)
	if e != nil {
		t.Fatal(e)
	}
	if result.(bool) == true {
		t.Fatal("express != nil fail")
	}
}

func TestEval_number(t *testing.T) {
	var aa = 2
	fmt.Println(Eval("2+1", Add, &aa, 1))
}

func BenchmarkEval(b *testing.B) {
	b.StopTimer()
	var aa = 1
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// var b,_=
		Eval("1/1", Ride, &aa, 1)
		//fmt.Println(b)
	}
}

type S struct {
	Name *string
}

func BenchmarkGetDeepValue(b *testing.B) {
	b.StopTimer()
	var a1 = int64(1)
	var v = reflect.ValueOf(&a1)
	var bb *int64
	bb = &a1
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GetDeepValue(v, bb)
	}
}

func TestEvalTakes(t *testing.T) {
	var strNode = ArgNode{
		value:     "a.B",
		values:    strings.Split("a.B", "."),
		valuesLen: 2,
	}
	var m = map[string]interface{}{"a": "B"}

	var r, _ = EvalTakes(strNode, m)
	fmt.Println(r)
}

type A struct {
	B  string
	C  C
	PC *C
}

type C struct {
	D  string
	PD *string
}

func TestEvalTakesStruct(t *testing.T) {
	var strNode = ArgNode{
		value:     "a.B",
		values:    strings.Split("a.B", "."),
		valuesLen: 2,
	}
	var m = A{B: "B", C: C{D: "D"}}

	var r, _ = EvalTakes(strNode, m)
	fmt.Println(r)
}

func TestEvalTakesStructD(t *testing.T) {
	var values = strings.Split("a.C.D", ".")
	var strNode = ArgNode{
		value:     "a.C.D",
		values:    values,
		valuesLen: len(values),
	}
	var m = A{B: "B", C: C{D: "D"}}

	var r, _ = EvalTakes(strNode, m)
	fmt.Println(r)
}

func TestEvalTakesStructPC(t *testing.T) {
	var values = strings.Split("a.PC.PD", ".")
	var strNode = ArgNode{
		value:     "a.PC",
		values:    values,
		valuesLen: len(values),
	}
	var pds = "pds"
	var c = C{
		D:  "D",
		PD: &pds,
	}
	var m = A{B: "B", C: C{D: "D"}, PC: &c}

	var r, e = EvalTakes(strNode, m)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r)
}

func BenchmarkEvalTakes(b *testing.B) {
	b.StopTimer()
	var values = strings.Split("a.PC.PD", ".")
	var strNode = ArgNode{
		value:     "a.PC.PD",
		values:    values,
		valuesLen: len(values),
	}
	var pds = "pds"
	var c = C{
		D:  "D",
		PD: &pds,
	}
	var m = A{B: "B", C: C{D: "D"}, PC: &c}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		EvalTakes(strNode, m)
	}
}
