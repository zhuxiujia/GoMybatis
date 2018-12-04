package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"testing"
)

func TestExpress(t *testing.T) {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}

	var evaluateParameters = make(map[string]interface{})

	evaluateParameters["activity"] = activity

	var expression = "activity.DeleteFlag == 1 or activity.DeleteFlag > 0 or activity.DeleteFlag ==0"
	evalExpression, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		panic(err)
	}
	result, err := evalExpression.Evaluate(evaluateParameters)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func Test_split(t *testing.T)  {
	var test = "a == 0 and a >= 0 or a < 0"
	var ns= GoMybatisSqlBuilder{}.split(&test)

	fmt.Println(test,ns)
}

func Test_bind_string(t *testing.T)  {
	var activity = example.Activity{
		Id:         "1",
		DeleteFlag: 1,
	}
	var evaluateParameters = make(map[string]interface{})
	evaluateParameters["activity"] = activity
	var expression = "'%' + activity.Id + '%'"
	evalExpression, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		panic(err)
	}
	result, err := evalExpression.Evaluate(evaluateParameters)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}