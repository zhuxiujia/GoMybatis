package GoMybatis

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhuxiujia/GoMybatis/example"
	"encoding/json"
	"log"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
)

var manager DefaultTransationManager

func TestManager(t *testing.T) {
	engine, err := Open("mysql", "*/test?charset=utf8&parseTime=True&loc=Local") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}
	var SessionFactory = SessionFactory{}.New(engine)
	var TransactionFactory = TransactionFactory{}.New(&SessionFactory)
	manager = DefaultTransationManager{}.New(&SessionFactory, &TransactionFactory)

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
	var OwnerId = utils.CreateUUID()
	var sql = "UPDATE `test`.`biz_property` SET `pool_amount`= (pool_amount+100) WHERE `id`='20180926172014a2a48a9491004603';"
	//todo proxy send error
	var dto = TransactionReqDTO{
		TransactionId: "1234",
		Status:        Transaction_Status_Pause,
		ActionType:    ActionType_Exec,
		Sql:           sql,
	}
	var result = manager.DoTransaction(manager, dto, OwnerId)
	fmt.Println(result.Exec)
	dto.Status = Transaction_Status_Commit
	manager.DoTransaction(manager, dto, OwnerId) //commit
	return nil
}

type TestPropertyServiceB struct{}

//单事务1
func (TestPropertyServiceB) Reduce(transactionId string, id string, amt int) error {
	var OwnerId = utils.CreateUUID()
	var sql = "UPDATE `test`.`biz_property` SET `pool_amount`= (pool_amount-100) WHERE `id`='20180926172014a2a48a9491004603';"
	//todo proxy send error
	var dto = TransactionReqDTO{
		TransactionId: "1234",
		Status:        Transaction_Status_Pause,
		ActionType:    ActionType_Exec,
		Sql:           sql,
	}
	var result = manager.DoTransaction(manager, dto, OwnerId)
	fmt.Println(result.Exec)
	dto.Status = Transaction_Status_Commit
	manager.DoTransaction(manager, dto, OwnerId) //commit
	return nil
}

type TestOrderService struct {
	TestPropertyServiceA TestPropertyServiceA //A微服务
	TestPropertyServiceB TestPropertyServiceB //B微服务
}

//嵌套事务
func (this TestOrderService) Transform(transactionId string, outid string, inId string, amount int) error {
	var OwnerId = utils.CreateUUID()
	transactionId = "2018092d6172014a2a4c8a949f1004623"
	var dto = TransactionReqDTO{
		TransactionId: transactionId,
		Status:        Transaction_Status_Pause,
		ActionType:    ActionType_Exec,
		Sql:           "",
	}
	manager.DoTransaction(manager, dto, OwnerId) //开启事务

	//事务id=2018092d6172014a2a4c8a949f1004623,已存在的事务不可提交commit，只能提交状态rollback和Pause
	var e1 = this.TestPropertyServiceB.Reduce(transactionId, outid, amount)
	if e1 != nil {
		return e1
	}
	//事务id=2018092d6172014a2a4c8a949f1004623,已存在的事务不可提交commit，只能提交状态rollback和Pause
	var e2 = this.TestPropertyServiceA.Add(transactionId, inId, amount)
	if e2 != nil {
		return e2
	}
	//事务id=2018092d6172014a2a4c8a949f1004623,原始事务可提交commit,rollback和Pause
	dto.Status = Transaction_Status_Commit
	manager.DoTransaction(manager, dto, OwnerId)
	return nil
}
