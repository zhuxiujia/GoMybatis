package GoMybatis

import (
	"testing"
	"github.com/zhuxiujia/GoMybatis/utils"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/example"
)

func TestTransationRM(t *testing.T) {
	var addr = "127.0.0.1:17235"
	go ServerTcp(addr, "mysql", example.MysqlUri) //事务服务器节点1

	var TransationRMClient = TransationRMClient{
		RetryTime: 3,
		Addr:      addr,
	}

	var result TransactionRspDTO

	var TransactionId = utils.CreateUUID() //服务站点事务

	TransationRMClient.Call(TransactionReqDTO{Status: Transaction_Status_Pause, TransactionId: TransactionId, Sql: "", ActionType: ActionType_Exec,}, &result)

	TransationRMClient.Call(TransactionReqDTO{Status: Transaction_Status_Rollback, TransactionId: TransactionId, Sql: "", ActionType: ActionType_Exec,}, &result)

	fmt.Println(result)

}
