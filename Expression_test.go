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

	var expression = "activity.Id == '2'"
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
