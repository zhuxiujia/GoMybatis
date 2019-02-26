package GoMybatis

import (
	"fmt"
	"testing"
)

func TestStringNode_Eval(t *testing.T) {
	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)
	var sNode = NodeString{
		value:      "#{name}",
		t:          NString,
		expressMap: []string{"name"},
		holder: &NodeConfigHolder{
			convert: &GoMybatisSqlArgTypeConvert{},
			proxy:   &proxy,
		},
	}

	var argMap = map[string]interface{}{
		"name": "sadf",
	}
	var r, e = sNode.Eval(argMap)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(r))
}

func BenchmarkStringNode_Eval(b *testing.B) {
	b.StopTimer()

	var proxy = ExpressionEngineProxy{}.New(&ExpressionEngineGoExpress{}, true)
	var sNode = NodeString{
		value:      "#{name}",
		t:          NString,
		expressMap: []string{"name"},
		holder: &NodeConfigHolder{
			convert: &GoMybatisSqlArgTypeConvert{},
			proxy:   &proxy,
		},
	}

	var argMap = map[string]interface{}{
		"name": "sadf",
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
