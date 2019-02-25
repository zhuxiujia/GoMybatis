package GoMybatis

import (
	"fmt"
	"testing"
)

func TestStringNode_Eval(t *testing.T) {
	var sNode = StringNode{
		value: "#{name}",
		t:     NString,
	}
	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)
	var r, e = sNode.Eval(map[string]interface{}{
		"SqlArgTypeConvert":      &GoMybatisSqlArgTypeConvert{},
		"*ExpressionEngineProxy": &proxy,
		"name":                   "sadf",
	})
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r)
}

func BenchmarkStringNode_Eval(b *testing.B) {
	b.StopTimer()
	var sNode = StringNode{
		value: "#{name}",
		t:     NString,
	}
	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var _, e = sNode.Eval(map[string]interface{}{
			"SqlArgTypeConvert":      &GoMybatisSqlArgTypeConvert{},
			"*ExpressionEngineProxy": &proxy,
			"name":                   "sadf",
		})
		if e != nil {
			b.Fatal(e)
		}
	}
}
