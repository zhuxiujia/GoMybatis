package GoMybatis

import (
	"log"
	"testing"
)

func TestLogStandard_Println(t *testing.T) {
	var stand = LogStandard{}
	//没有设置func，使用系统log打印
	stand.Println([]byte("hello"))

	//设置func，使用func打印
	stand.PrintlnFunc = func(v []byte) {
		log.Println(string(v), "println on PrintlnFunc()")
	}

	stand.Println([]byte("hello"))
}
