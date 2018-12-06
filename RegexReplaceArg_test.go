package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
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

	var param map[string]SqlArg
	param = make(map[string]SqlArg)
	var age = 11
	fmt.Println("age=", age)
	var bean = TestBean{
		Name: "father",
		Child: TestBeanChild{
			Name: "child",
			Age:  &age,
		},
	}
	param["bean"] = SqlArg{
		Value: bean,
		Type:  reflect.TypeOf(bean),
	}
	defer utils.CountMethodUseTime(time.Now(), "Test_Access_Arg", time.Millisecond)
	var string = "-----#{bean.Name}------#{bean.Child.Age}---"

	for i := 0; i < 1; i++ {
		var arg = replaceArg(string, param, GoMybatisSqlArgTypeConvert{})
		fmt.Println(arg)
	}
}
