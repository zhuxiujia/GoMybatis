package GoMybatis

import (
	"GoMybatis/utils"
	"bytes"
)

type LogSystem struct {
	log     Log
	logChan chan []byte
	started bool
}

//logImpl:日志实现类,queueLen:消息队列缓冲长度
func (it LogSystem) New(logImpl Log, queueLen int) (LogSystem, error) {
	if it.started == true {
		return it, utils.NewError("LogSystem", "log system is started!")
	}
	if logImpl == nil {
		logImpl = &LogStandard{}
	}
	it.logChan = make(chan []byte, queueLen)
	it.log = logImpl
	//启动接受者
	go it.receiver()
	it.started = true
	return it, nil
}

//关闭日志系统和队列
func (it *LogSystem) Close() error {
	close(it.logChan)
	it.started = false
	return nil
}

//日志发送者
func (it *LogSystem) SendLog(logs ...string) error {
	if it.started == false {
		return utils.NewError("LogSystem", "no log Receiver! you must call go GoMybatis.LogSystem{}.New()")
	}
	var buf bytes.Buffer
	for _, v := range logs {
		buf.WriteString(v)
	}
	it.logChan <- buf.Bytes()
	return nil
}

//日志接受者
func (it *LogSystem) receiver() error {
	for {
		logs := <-it.logChan
		it.log.Println(logs)
	}
	return nil
}
