package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
)

type TestBean struct {
	Name  string
	Child TestBeanChild
}
type TestBeanChild struct {
	Name string
	Age  *int
}

func Test_Access_Arg(t *testing.T) {

	var param map[string]interface{}
	param = make(map[string]interface{})
	var age = 11
	fmt.Println("age=", age)
	var bean = TestBean{
		Name: "father",
		Child: TestBeanChild{
			Name: "child",
			Age:  &age,
		},
	}
	param["bean"] = bean
	defer utils.CountMethodUseTime(time.Now(), "Test_Access_Arg", time.Millisecond)
	var string = "-----#{bean.Name}------#{bean.Child.Age}---"

	var arg, err = replaceArg(string, param, GoMybatisSqlArgTypeConvert{}, &ExpressionEngineJee{})
	if err != nil {
		t.Fatal(err)
	}
	if arg != "-----father------11---" {
		t.Fatal("replaceArgFail")
	}
	fmt.Println(arg)
}
