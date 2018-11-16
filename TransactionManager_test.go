package GoMybatis

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis/example"
	"encoding/json"
	"log"
	"fmt"
)

func TestManager(t *testing.T) {
	engine, err := Open("mysql", "*/test?charset=utf8&parseTime=True&loc=Local") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}
	var SessionFactory = SessionFactory{}.New(engine)
	var TransactionFactory = TransactionFactory{}.New(&SessionFactory)
	var manager = DefaultTransationManager{}.New(&SessionFactory, &TransactionFactory)

	var dto = TransactionReqDTO{
		TransactionId: "1234",
		Status:        Transaction_Status_Pause,
		ActionType:    ActionType_Query,
		Sql:           "select * from biz_activity where delete_flag = 1",
	}

	//start
	var result = manager.DoTransaction(manager, dto)

	printData(result)

	dto.Sql = "UPDATE `test`.`biz_activity` SET `name`='rs-updated' WHERE `id`='159'"
	dto.Status = Transaction_Status_Pause
	dto.ActionType = ActionType_Exec
	manager.DoTransaction(manager, dto)

	dto.Status = Transaction_Status_Rollback
	manager.DoTransaction(manager, dto)

}

func printData(result TransactionRspDTO) {
	if result.Error != "" {
		log.Println(result.Error)
		return
	}
	var Activity []example.Activity
	Unmarshal(result.Query, &Activity)
	var b, _ = json.Marshal(Activity)
	log.Println("Activity Json=", string(b))
}

//测试案例
type TestUser struct {
	Id     string
	Amount int
}

type TestPropertyServiceA struct{}

//单事务2
func (TestPropertyServiceA) Add(transactionId string, id string, amt int) error {
	var sql="UPDATE `test`.`biz_property` SET `pool_amount`= (pool_amount+100) WHERE `id`='20180926172014a2a48a9491004603';"
	fmt.Println(sql)
	//todo proxy send error
	return nil
}

type TestPropertyServiceB struct{}

//单事务1
func (TestPropertyServiceB) Reduce(transactionId string, id string, amt int) error {
	var sql="UPDATE `test`.`biz_property` SET `pool_amount`= (pool_amount-100) WHERE `id`='20180926172014a2a48a9491004603';"
	fmt.Println(sql)
	//todo proxy send error
	return nil
}

type TestOrderService struct {
	TestPropertyServiceA TestPropertyServiceA
	TestPropertyServiceB TestPropertyServiceB
}

//嵌套事务
func (this TestOrderService) Transform(transactionId string, outid string, inId string, amount int) error {
	//事务id=1234
	var e1 = this.TestPropertyServiceB.Reduce(transactionId, outid, amount)
	if e1 != nil {
		return e1
	}
	var e2 = this.TestPropertyServiceA.Add(transactionId, inId, amount)
	if e2 != nil {
		return e2
	}
	//todo proxy send error
	return nil
}
