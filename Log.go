package GoMybatis

type Log interface {
	QueueLen() int//日志消息队列长度
	Println(messages []byte)//日志输出方法实现
}
