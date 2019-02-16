package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
)

func TestExpress(t *testing.T) {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}
	var engine = ExpressionEngineGovaluate{}
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

func TestExpressionEngineGovaluateTakeValue(t *testing.T) {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}
	var engine = ExpressionEngineGovaluate{}
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

func BenchmarkExpress(b *testing.B) {
	b.StopTimer()
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var engine = ExpressionEngineGovaluate{}
	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = &activity
	var expression = "activity.DeleteFlag == 1 and activity.DeleteFlag != 0 "
	evalExpression, err := engine.Lexer(expression)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Eval(evalExpression, evaluateParameters, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestTpsExpressionEngineGovaluate(t *testing.T) {

	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var engine = ExpressionEngineGovaluate{}
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


func TestExpress_nil(t *testing.T) {
	var engine = ExpressionEngineGovaluate{}
	var evaluateParameters = make(map[string]interface{})

	var namePtr *string

	evaluateParameters["name"] = namePtr



	var expression = "name != nil"
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
