package Knetic

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/Knetic/govaluate"
	"reflect"

	"testing"
)

func Test_evaluate(t *testing.T) {

	var a= 1
	var param = make(map[string]interface{})
	param["a"] = &a

	fmt.Println(reflect.ValueOf(param["a"]).Elem().Interface())

	evalExpression, err := govaluate.NewEvaluableExpression( "*a == 1")
	if err != nil {
		t.Error(err)
	}
	result, err := evalExpression.Evaluate(param)
	fmt.Println(err)
	fmt.Println(result)
}
