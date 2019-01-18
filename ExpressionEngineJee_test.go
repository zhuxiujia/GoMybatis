package GoMybatis

import (
	"fmt"
	"testing"
)

func TestExpressionEngineJee_Eval(t *testing.T) {
	var m = make(map[string]interface{})
	m["a"] = nil
	var engine = ExpressionEngineJee{}
	var result, error = engine.Eval(".a == null", m)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println(result)
}

func TestExpressionEngineJee_Eval_add(t *testing.T) {
	var m = make(map[string]interface{})
	m["a"] = 1
	var engine = ExpressionEngineJee{}
	var result, error = engine.Eval(".a + 1", m)
	if error != nil {
		t.Fatal(error)
	}
	fmt.Println(result)
}

func BenchmarkExpressionEngineJee_Eval(b *testing.B) {
	b.StopTimer()
	var m = make(map[string]interface{})
	m["a"] = nil
	var engine = ExpressionEngineJee{}

	//start
	b.StartTimer()
	for i := 0; i < b.N; i++ {

		var result, error = engine.Eval(".a == null", m)
		if error != nil {
			b.Fatal(error)
		}
		if result == nil {
			b.Fatal("eval fail")
		}
	}
}
