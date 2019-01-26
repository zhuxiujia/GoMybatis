package GoMybatis

import "log"

type LogStandard struct {
	PrintlnFunc func(messages []byte) //日志输出方法实现
}

//日志消息队列长度
func (it *LogStandard) QueueLen() int {
	//默认50万个日志消息缓存队列
	return 500000
}

//日志输出方法实现
func (it *LogStandard) Println(v []byte) {
	if it.PrintlnFunc != nil {
		it.PrintlnFunc(v)
	} else {
		log.Println(string(v))
	}
}
