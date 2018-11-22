package GoMybatis

import (
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
	"testing"
)

func TestTransationRM(t *testing.T) {
	var addr = "127.0.0.1:17235"
	go ServerTcp(addr, "mysql", example.MysqlUri) //事务服务器节点1

	var TransationRMClient = TransationRMClient{
		RetryTime: 3,
		Addr:      addr,
	}

	var transationRMServerSession = TransationRMSession{
		Client:  &TransationRMClient,
		OwnerId: "1234",
	}

	var e error

	e = transationRMServerSession.Begin()
	if e != nil {
		panic(e)
	}

	result, e := transationRMServerSession.Exec("UPDATE `test`.`biz_activity` SET `name`='rs168-10' WHERE `id`='170';")
	if e != nil {
		panic(e)
	}

	fmt.Println(result)

	e = transationRMServerSession.Commit()
	if e != nil {
		panic(e)
	}
}
