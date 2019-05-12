package GoMybatis

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const Adapter_FormateDate = `2006-01-02 15:04:05`

//Sql内容类型转换器
type GoMybatisSqlArgTypeConvert struct {
}

//Sql内容类型转换器
func (it GoMybatisSqlArgTypeConvert) Convert(argValue interface{}, argType reflect.Type) string {
	if argType == nil {
		argType = reflect.TypeOf(argValue)
	}
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
	case reflect.String:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(argValue.(string))
		argStr.WriteString(`'`)
		return argStr.String()
	case reflect.Struct:
		if strings.Contains(argType.String(), "time.Time") {
			var argStr bytes.Buffer
			argStr.WriteString(`'`)
			argStr.WriteString(argValue.(time.Time).Format(Adapter_FormateDate))
			argStr.WriteString(`'`)
			return argStr.String()
		}
		break
	}
	return it.toString(argValue, argType)
}

func (it GoMybatisSqlArgTypeConvert) toString(value interface{}, argType reflect.Type) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
