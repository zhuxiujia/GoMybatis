package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/engines"
	"github.com/zhuxiujia/GoMybatis/example"
	"fmt"
	"testing"
)

func TestExpressionEngineProxy_Eval(t *testing.T) {
	var engine = ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, false)
	var lexer, err = engine.Lexer("foo")
	if err != nil {
		t.Fatal(err)
	}
	var arg = make(map[string]interface{})
	arg["foo"] = "Bar"
	result, err := engine.Eval(lexer, arg, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.(string) != "Bar" {
		t.Fatal("result != 'Bar'")
	}
	fmt.Println(result)
}

func TestExpressionEngineProxy_Lexer(t *testing.T) {
	var engine = ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, false)
	var _, err = engine.Lexer("foo")
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkExpressionEngineProxy_Eval(b *testing.B) {
	b.StopTimer()
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var engine = ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, false)
	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = &activity

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var expression = "activity.DeleteFlag == 1 and activity.DeleteFlag != 0 "
		evalExpression, err := engine.Lexer(expression)
		if err != nil {
			b.Fatal(err)
		}
		result, err := engine.Eval(evalExpression, evaluateParameters, 0)
		if err != nil {
			b.Fatal(err)
		}
		if result.(bool) {

		}
	}
}

func BenchmarkExpressionEngineProxy_Eval_each(b *testing.B) {
	b.StopTimer()
	var engine = ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true)
	var evaluateParameters = make(map[string]interface{})
	var name = "dsafas"
	evaluateParameters["activity"] = &name

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for run := 0; run < 8; run++ {
			var expression = "activity"
			evalExpression, err := engine.Lexer(expression)
			if err != nil {
				b.Fatal(err)
			}
			_, err = engine.Eval(evalExpression, evaluateParameters, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkExpressionEngineProxy_LexerAndEval(b *testing.B) {
	b.StopTimer()
	var engine = ExpressionEngineProxy{}.New(&engines.ExpressionEngineGoExpress{}, true)
	var evaluateParameters = make(map[string]interface{})
	var name = "dsafas"
	evaluateParameters["activity"] = name
	//evaluateParameters["func_activity != nil"] = func(arg map[string]interface{}) interface{} {
	//	return arg["activity"] != nil
	//}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for run := 0; run < 8; run++ {
			var expression = "activity != nil"
			_, err := engine.LexerAndEval(expression, evaluateParameters)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
