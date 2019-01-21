package GoMybatis

import (
	"fmt"
	"log"
	"testing"
)

func TestLogSystem_SendLog(t *testing.T) {
	var endChan = make(chan int)
	var stdLog = LogStandard{
		PrintlnFunc: func(v []byte) {
			log.Println(string(v))
			endChan <- 1
		},
	}
	var system, err = LogSystem{}.New(&stdLog, 1000000, Log_Mode_async)
	if err != nil {
		t.Fatal(err)
	}
	err = system.SendLog("hello")
	if err != nil {
		t.Fatal(err)
	}
	var data = <-endChan
	fmt.Println(data)
}

func TestLogSystem_Close(t *testing.T) {

	var endChan = make(chan int)
	var stdLog = LogStandard{
		PrintlnFunc: func(v []byte) {
			log.Println(string(v))
			endChan <- 1
		},
	}
	var system, err = LogSystem{}.New(&stdLog, 1000000, Log_Mode_async)
	if err != nil {
		t.Fatal(err)
	}
	err = system.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = system.SendLog("hello")
	if err == nil {
		t.Fatal(err)
	}
	fmt.Println(err)
}
