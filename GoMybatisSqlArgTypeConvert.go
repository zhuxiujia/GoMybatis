package GoMybatis

import (
	"bytes"
	"fmt"
	"time"
)

const Adapter_FormateDate = `2006-01-02 15:04:05`

//Sql内容类型转换器
type GoMybatisSqlArgTypeConvert struct {
}

//Sql内容类型转换器
func (it GoMybatisSqlArgTypeConvert) Convert(argValue interface{}) string {
	if argValue == nil {
		return "''"
	}
	switch argValue.(type) {
	case string:
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		//argStr.WriteString(argValue.(string))
		argStr.WriteString(antiSqlInjectionExp(argValue.(string)))
		argStr.WriteString(`'`)
		return argStr.String()
	case *string:
		var v = argValue.(*string)
		if v == nil {
			return "''"
		}
		var argStr bytes.Buffer
		argStr.WriteString(`'`)
		//argStr.WriteString(*v)
		argStr.WriteString(antiSqlInjectionExp(*v))
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

	case int, int16, int32, int64, float32, float64:
		return fmt.Sprint(argValue)
	case *int:
		var v = argValue.(*int)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *int16:
		var v = argValue.(*int16)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *int32:
		var v = argValue.(*int32)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *int64:
		var v = argValue.(*int64)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *float32:
		var v = argValue.(*float32)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	case *float64:
		var v = argValue.(*float64)
		if v == nil {
			return ""
		}
		return fmt.Sprint(*v)
	}

	return it.toString(argValue)
}

func (it GoMybatisSqlArgTypeConvert) toString(argValue interface{}) string {
	if argValue == nil {
		return ""
	}
	return fmt.Sprint(argValue)
}
// 字符串防Sql注入[将字符'替换为\']
// 不足之处：
// １）破坏了原值；
// ２）本框架与其它使用别的框架（如java的mybatis）等共用数据源时，可能有到导致处理结果不
func antiSqlInjectionStringExp(str string) string{
	return strings.ReplaceAll(str,"'","\'")
}