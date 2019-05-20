package GoMybatis

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

const Adapter_FormateDate = `2006-01-02 15:04:05`

//Sql内容类型转换器
type GoMybatisSqlArgTypeConvert struct {
}

//Sql内容类型转换器
func (it GoMybatisSqlArgTypeConvert) Convert(argValue interface{}) string {
	//if argType == nil {
	//	argType = reflect.TypeOf(argValue)
	//}
	if argValue == nil {
		return "''"
	}
	var argValueV = reflect.ValueOf(argValue)
	if !argValueV.IsValid() {
		return "''"
	}
	switch argValue.(type) {
	case string:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(argValue.(string))
		argStr.WriteString(`'`)
		return argStr.String()
	case *string:
		var v = argValue.(*string)
		if v == nil {
			return "''"
		}
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(*v)
		argStr.WriteString(`'`)
		return argStr.String()
	case bool:
		if argValue.(bool) {
			return "true"
		} else {
			return "false"
		}
	case *bool:
		var v = argValue.(*bool)
		if v == nil {
			return "''"
		}
		if *v {
			return "true"
		} else {
			return "false"
		}
	case time.Time:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(argValue.(time.Time).Format(Adapter_FormateDate))
		argStr.WriteString(`'`)
		return argStr.String()
	case *time.Time:
		var timePtr = argValue.(*time.Time)
		if timePtr == nil {
			return "''"
		}
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		argStr.WriteString(timePtr.Format(Adapter_FormateDate))
		argStr.WriteString(`'`)
		return argStr.String()

	}

	return it.toString(argValue)
}

func (it GoMybatisSqlArgTypeConvert) toString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
