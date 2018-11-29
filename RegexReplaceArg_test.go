package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
	"reflect"
)

type TestBean struct {
	Name  string
	Child TestBeanChild
}
type TestBeanChild struct {
	Name string
	Age  int
}

func Test_Access_Arg(t *testing.T) {

	var param map[string]SqlArg
	param = make(map[string]SqlArg)

	var bean=TestBean{
		Name: "father",
		Child: TestBeanChild{
			Name: "child",
			Age:  11,
		},
	}
	param["bean"] = SqlArg{
		Value:bean,
		Type:reflect.TypeOf(bean),
	}
	defer utils.CountMethodUseTime(time.Now(), "Test_Access_Arg", time.Millisecond)
	var string = "-----#{bean.Name}------#{bean.Child.age}---"

	for i := 0; i < 1; i++ {
		var arg = replaceArg(string, param, GoMybatisSqlArgTypeConvert{})
		fmt.Println(arg)
	}
}
