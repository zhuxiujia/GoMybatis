package GoMybatis

import (
	"testing"
	"github.com/zhuxiujia/GoMybatis/utils"
	"fmt"
)

func TestTransationRM(t *testing.T) {
	var addr="127.0.0.1:17235"
	go ServerTransationRM(addr, "mysql", "*?charset=utf8&parseTime=True&loc=Local")

	var TransationRMClient = TransationRMClient{
		RetryTime:3,
		Addr:addr,
	}

	var result TransactionRspDTO

	TransationRMClient.Call(TransactionReqDTO{
		Status:        Transaction_Status_Pause,
		TransactionId: utils.CreateUUID(),
		Sql:           "",
		ActionType:    ActionType_Exec,
	}, &result)

	TransationRMClient.Call(TransactionReqDTO{
		Status:        Transaction_Status_Rollback,
		TransactionId: utils.CreateUUID(),
		Sql:           "",
		ActionType:    ActionType_Exec,
	}, &result)

	fmt.Println(result)


}
