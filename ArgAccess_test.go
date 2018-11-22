package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"github.com/zhuxiujia/GoMybatis/utils"
	"testing"
	"time"
)

func Test_Access_Arg(t *testing.T) {

	var param map[string]interface{}
	param = make(map[string]interface{})
	param["activity"] = example.Activity{
		Name: "aaaa",
	}

	defer utils.CountMethodUseTime(time.Now(), "Test_Access_Arg", time.Millisecond)
	var string = "---#{activity.name}------#{activity.id}---"
	fmt.Println("start")

	for i := 0; i < 1; i++ {
		var arg = replaceArg(string, param, DefaultSqlTypeConvertFunc)
		fmt.Println(arg)
	}
}
