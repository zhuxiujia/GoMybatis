package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/utils"
	"reflect"
	"time"
)

const Adapter_FormateDate = `2006-01-02 15:04:05`

//Sql内容类型转换器
type GoMybatisSqlArgTypeConvert struct {
}

//Sql内容类型转换器
func (it GoMybatisSqlArgTypeConvert) Convert(arg SqlArg) string {
	var argValue = arg.Value
	var argType = arg.Type
	if argValue == nil {
		return "''"
	}
	switch argType.Kind() {
	case reflect.Bool:
		if argValue.(bool) {
			return "true"
		} else {
			return "false"
		}
		break
	case reflect.String:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(it.toString(&arg))
		argStr.WriteString(`'`)
		return argStr.String()
	case reflect.Struct:
		if argType.String() == Adapter_DateType {
			var argStr bytes.Buffer
			argStr.WriteString(`'`)
			argStr.WriteString(argValue.(time.Time).Format(Adapter_FormateDate))
			argStr.WriteString(`'`)
			return argStr.String()
		}
		break
	}
	return it.toString(&arg)
}

func (it GoMybatisSqlArgTypeConvert) toString(value *SqlArg) string {
	if value.Value == nil {
		return ""
	}
	//return fmt.Sprint(value.Value)
	return utils.GetValue(value.Value, value.Type)
}
