package GoMybatis

import (
	"fmt"
	"log"
	"testing"
)

func TestLogStandard_Println(t *testing.T) {
	var stand = LogStandard{}
	//没有设置func，使用系统log打印
	stand.Println("hello")

	//设置func，使用func打印
	stand.PrintlnFunc = func(v ...string) {
		log.Println(fmt.Sprint(v))
	}

	stand.Println("hello")
}
