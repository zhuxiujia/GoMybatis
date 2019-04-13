package engines

import (
	"GoMybatis/example"
	"GoMybatis/utils"
	"fmt"
	"testing"
	"time"
)

func TestExpressionEngineExpr_Eval(t *testing.T) {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}
	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = &activity

	var expression = "activity.DeleteFlag == 1 || activity.DeleteFlag > 0 "
	evalExpression, err := engine.Lexer(expression)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Eval(evalExpression, evaluateParameters, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

type TestPtr struct {
	Age *int
}

func TestExpressionEngineExpr_Struct(t *testing.T) {
	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})
	//var p=1
	var age = 1
	evaluateParameters["obj"] = &TestPtr{
		Age: &age,
	}
	var expression = "obj.Age"
	evalExpression, err := engine.Lexer(expression)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Eval(evalExpression, evaluateParameters, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestExpressionEngineExpr_Nil_Null(t *testing.T) {
	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})
	//var p=1
	evaluateParameters["startTime"] = nil
	var nmap = makeArgInterfaceMap(evaluateParameters)
	var expression = "startTime == nil"
	evalExpression, err := engine.Lexer(expression)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Eval(evalExpression, nmap, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func BenchmarkExpressionEngineExprNil_Null(b *testing.B) {
	b.StopTimer()
	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})
	//var p=1
	evaluateParameters["startTime"] = nil
	var nmap = makeArgInterfaceMap(evaluateParameters)
	evalExpression, err := engine.Lexer("startTime == nil")
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Eval(evalExpression, nmap, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExpressionEngineExprNilTakeValue(b *testing.B) {
	b.StopTimer()
	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})
	var startTime *string
	var startTimeV *string

	var s = "12345"
	startTimeV = &s
	evaluateParameters["startTime"] = startTime
	evaluateParameters["startTimeValue"] = startTimeV
	var nmap = makeArgInterfaceMap(evaluateParameters)
	evalExpression, err := engine.Lexer("startTime == nil")
	if err != nil {
		b.Fatal(err)
	}
	takeValueExpression, err := engine.Lexer("startTimeValue")
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < 8; k++ {
			_, err := engine.Eval(evalExpression, nmap, 0)
			if err != nil {
				b.Fatal(err)
			}
			_, err = engine.Eval(takeValueExpression, nmap, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func makeArgInterfaceMap(args map[string]interface{}) map[string]interface{} {
	var m = make(map[string]interface{})
	if args != nil {
		for k, v := range args {
			m[k] = v
		}
	}
	return m
}

func TestExpressionEngineExpr_Eval_TakeValue(t *testing.T) {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}
	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = &activity

	var expression = "activity.DeleteFlag"
	evalExpression, err := engine.Lexer(expression)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Eval(evalExpression, evaluateParameters, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func BenchmarkExpressionEngineExpr_Eval(b *testing.B) {
	b.StopTimer()
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = &activity
	var expression = "activity.DeleteFlag == 1 and activity.DeleteFlag != 0 "
	evalExpression, err := engine.Lexer(expression)
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = engine.Eval(evalExpression, evaluateParameters, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestTpsExpressionEngineExpr_Eval(t *testing.T) {

	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var engine = ExpressionEngineExpr{}
	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = &activity

	defer utils.CountMethodTps(10000, time.Now(), "ExpressionEngineGovaluate")
	for i := 0; i < 10000; i++ {
		var expression = "activity.DeleteFlag == 1 || activity.DeleteFlag > 0 "
		evalExpression, err := engine.Lexer(expression)
		if err != nil {
			t.Fatal(err)
		}
		result, err := engine.Eval(evalExpression, evaluateParameters, 0)
		if err != nil {
			t.Fatal(err)
		}
		if result.(bool) {

		}
	}
}
