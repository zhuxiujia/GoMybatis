package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
)

type LogSystem struct {
	log     Log
	logChan chan string
	started bool
}

func (this LogSystem) New(l Log) (LogSystem, error) {
	if this.started == true {
		return this, utils.NewError("LogSystem", "log system is started!")
	}
	if l == nil {
		l = &LogStandard{}
	}
	this.logChan = make(chan string)
	this.log = l
	go this.receiver()
	this.started = true
	return this, nil
}

func (this *LogSystem) Close() ( error) {
	close(this.logChan)
	this.started = false
	return nil
}

func (this *LogSystem) SendLog(logs string) error {
	if this.started == false {
		return utils.NewError("LogSystem", "no log Receiver! you must call go GoMybatis.LogSystem{}.New()")
	}
	this.logChan <- logs
	return nil
}

func (this *LogSystem) receiver() error {
	logs := <-this.logChan
	this.log.Println(logs)
	return nil
}
