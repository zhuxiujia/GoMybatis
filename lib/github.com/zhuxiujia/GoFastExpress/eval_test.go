package GoFastExpress

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	fmt.Println(Eval("", nil, "a"))
}

func TestEvalNil(t *testing.T) {
	var b *string
	var result, e = Eval("!=", nil, b)
	if e != nil {
		t.Fatal(e)
	}
	if result.(bool) == true {
		t.Fatal("express != nil fail")
	}
}

func TestEval_number(t *testing.T) {
	var aa = 2
	fmt.Println(Eval(Add, &aa, 1))
}

func BenchmarkEval(b *testing.B) {
	b.StopTimer()
	var aa = 1
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// var b,_=
		Eval(Ride, &aa, 1)
		//fmt.Println(b)
	}
}

type S struct {
	Name *string
}

func TestEvalTake(t *testing.T) {
	var aa = 2
	fmt.Println(EvalTake("", &aa))

	var s = S{
		Name: nil,
	}
	fmt.Println(EvalTake("s.Name", &s))
}

func BenchmarkEvalTake(b *testing.B) {
	b.StopTimer()
	//var aa=1
	var s = S{
		Name: nil,
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// var b=
		//EvalTake("aa", aa)
		EvalTake("s.Name", &s)
		//fmt.Println(b)
	}
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
