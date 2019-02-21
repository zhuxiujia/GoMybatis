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

	var engine = ExpressionEngineProxy{}.New(&ExpressionEngineJee{}, true)
	var arg, err = replaceArg(string, param, GoMybatisSqlArgTypeConvert{}, &engine)
	if err != nil {
		t.Fatal(err)
	}
	if arg != "-----'father'------11---" {
		t.Fatal("replaceArgFail", "arg=", arg)
	}
	fmt.Println(arg)
}

func BenchmarkSplite(b *testing.B) {
	b.StopTimer()
	var str = "#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//sqlArgRegex.FindAllString(str,-1)
		//strings.Split(str,"#{")
		//strings.SplitAfter()
		FindAllExpressConvertString(str)
	}
}

func TestFindAllExpressConvertString(t *testing.T) {
	var str = "#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}"
	fmt.Println(FindAllExpressConvertString(str))
}
