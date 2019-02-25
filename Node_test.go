package GoMybatis

import (
	"fmt"
	"testing"
)

func TestStringNode_Eval(t *testing.T) {
	var sNode = NodeString{
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
	var sNode = NodeString{
		value: "#{name}",
		t:     NString,
	}
	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)

	var argMap = map[string]interface{}{
		"SqlArgTypeConvert":      &GoMybatisSqlArgTypeConvert{},
		"*ExpressionEngineProxy": &proxy,
		"name":                   "sadf",
	}
	//var lex,e=proxy.Lexer("name")
	//if e!=nil{
	//	b.Fatal(e)
	//}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//if true{
		//   proxy.Eval(lex,argMap,0)
		//	continue
		//}
		var _, e = sNode.Eval(argMap)
		if e != nil {
			b.Fatal(e)
		}
	}
}
