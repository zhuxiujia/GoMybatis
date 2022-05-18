package GoMybatis

import "log"

type LogStandard struct {
	PrintlnFunc func(messages ...string) //日志输出方法实现
}

//日志输出方法实现
func (it *LogStandard) Println(v ...string) {
	if it.PrintlnFunc != nil {
		it.PrintlnFunc(v...)
	} else {
		log.Println(v)
	}
}
