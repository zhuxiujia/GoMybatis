package GoMybatis

import (
	"bytes"
	"github.com/zhuxiujia/GoMybatis/utils"
	"time"
)

type LogMode int

const (
	Log_Mode_async     LogMode = iota //异步，非实时
	Log_Mode_Real_Time                //同步,实时
)

type LogSystem struct {
	log     Log
	logChan chan []byte
	started bool
	mode    LogMode
}

//logImpl:日志实现类,queueLen:消息队列缓冲长度
func (this LogSystem) New(logImpl Log, queueLen int, mode LogMode) (LogSystem, error) {
	if this.started == true {
		return this, utils.NewError("LogSystem", "log system is started!")
	}
	if logImpl == nil {
		logImpl = &LogStandard{}
	}
	this.logChan = make(chan []byte, queueLen)
	this.log = logImpl
	this.mode = mode

	//启动接受者
	go this.receiver()
	this.started = true
	return this, nil
}

//关闭日志系统和队列
func (this *LogSystem) Close() error {
	close(this.logChan)
	this.started = false
	return nil
}

//日志发送者
func (this *LogSystem) SendLog(logs ...string) error {
	if this.started == false {
		return utils.NewError("LogSystem", "no log Receiver! you must call go GoMybatis.LogSystem{}.New()")
	}
	var buf bytes.Buffer
	for _, v := range logs {
		buf.WriteString(v)
	}
	this.logChan <- buf.Bytes()
	return nil
}

//日志接受者
func (this *LogSystem) receiver() error {
	for {
		if this.mode == Log_Mode_async {
			//非实时，加延迟，出让cpu时间片，否则日志占用cpu过多将拖慢系统
			time.Sleep(50 * time.Millisecond)
		} else if this.mode == Log_Mode_Real_Time {
			//实时，直接执行操作
		}
		logs := <-this.logChan
		this.log.Println(logs)
	}
	return nil
}
