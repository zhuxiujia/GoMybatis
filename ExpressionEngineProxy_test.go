package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"testing"
)

func TestExpressionEngineProxy_Eval(t *testing.T) {
	var engine=ExpressionEngineProxy{}.New(&ExpressionEngineExpr{},false)
	var lexer,err=engine.Lexer("foo")
	if err!= nil{
		t.Fatal(err)
	}
	var arg=make(map[string]interface{})
	arg["foo"]="Bar"
	result,err:=engine.Eval(lexer,arg,0)
	if err!= nil{
		t.Fatal(err)
	}
	if result.(string) != "Bar"{
		t.Fatal("result != 'Bar'")
	}
	fmt.Println(result)
}

func TestExpressionEngineProxy_Lexer(t *testing.T) {
	var engine=ExpressionEngineProxy{}.New(&ExpressionEngineExpr{},false)
	var _,err=engine.Lexer("foo")
	if err!= nil{
		t.Fatal(err)
	}
}

func BenchmarkExpressionEngineProxy_Eval(b *testing.B) {
	b.StopTimer()
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var engine = ExpressionEngineProxy{}.New(&ExpressionEngineExpr{},false)
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