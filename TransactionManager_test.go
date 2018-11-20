package GoMybatis

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/zhuxiujia/GoMybatis/utils"
	"strconv"
	"github.com/zhuxiujia/GoMybatis/example"
)

var manager DefaultTransationManager

func TestManager(t *testing.T) {
	engine, err := Open(example.MysqlDriverName, example.MysqlUri) //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}
	var SessionFactory = SessionFactory{}.New(engine)
	var TransactionFactory = TransactionFactory{}.New(&SessionFactory)
	manager = DefaultTransationManager{}.New(&SessionFactory, &TransactionFactory)

	var TestPropertyServiceA TestPropertyServiceA
	var TestPropertyServiceB TestPropertyServiceB
	var TestOrderService = TestOrderService{
		TestPropertyServiceA: TestPropertyServiceA,
		TestPropertyServiceB: TestPropertyServiceB,
	}
	TestOrderService.Transform(utils.CreateUUID(),  "20181023162632152fd236d6877ff4","20180926172013b85403d3715d46ed", 100)
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
	var sql = "UPDATE `test`.`biz_property` SET `pool_amount`= (pool_amount+" + strconv.Itoa(amt) + ") WHERE `user_id`='" + id + "';"
	//todo proxy send error
	var dto = TransactionReqDTO{
		TransactionId: transactionId,
		OwnerId:       OwnerId,
		Status:        Transaction_Status_Prepare,
		ActionType:    ActionType_Exec,
		Sql:           sql,
	}
	var result = manager.DoTransaction(dto)
	fmt.Println(result.Exec)
	dto.Status = Transaction_Status_Commit
	manager.DoTransaction(dto) //commit
	return nil
}

type TestPropertyServiceB struct{}

//单事务1
func (TestPropertyServiceB) Reduce(transactionId string, id string, amt int) error {
	var OwnerId = utils.CreateUUID()
	var sql = "UPDATE `test`.`biz_property` SET `pool_amount`= (pool_amount-" + strconv.Itoa(amt) + ") WHERE `user_id`='" + id + "';"
	//todo proxy send error
	var dto = TransactionReqDTO{
		TransactionId: transactionId,
		OwnerId:       OwnerId,
		Status:        Transaction_Status_Prepare,
		ActionType:    ActionType_Exec,
		Sql:           sql,
	}
	var result = manager.DoTransaction(dto)
	fmt.Println(result.Exec)
	dto.Status = Transaction_Status_Commit
	manager.DoTransaction(dto) //commit
	return nil
}

type TestOrderService struct {
	TestPropertyServiceA TestPropertyServiceA //A微服务
	TestPropertyServiceB TestPropertyServiceB //B微服务
}

//嵌套事务
func (this TestOrderService) Transform(transactionId string, outid string, inId string, amount int) error {
	var OwnerId = utils.CreateUUID()
	var dto = TransactionReqDTO{
		TransactionId: transactionId,
		OwnerId:       OwnerId,
		Status:        Transaction_Status_Prepare,
		ActionType:    ActionType_Exec,
		Sql:           "",
	}
	manager.DoTransaction(dto) //开启事务

	//事务id=2018092d6172014a2a4c8a949f1004623,已存在的事务不可提交commit，只能提交状态rollback和Pause
	var e1 = this.TestPropertyServiceB.Reduce(transactionId, outid, amount)
	if e1 != nil {
		return e1
	}

	dto.Status=Transaction_Status_Rollback
	manager.DoTransaction(dto)

	//事务id=2018092d6172014a2a4c8a949f1004623,已存在的事务不可提交commit，只能提交状态rollback和Pause
	var e2 = this.TestPropertyServiceA.Add(transactionId, inId, amount)
	if e2 != nil {
		return e2
	}



	//manager.Rollback(transactionId)
	//事务id=2018092d6172014a2a4c8a949f1004623,原始事务可提交commit,rollback和Pause
	dto.Status = Transaction_Status_Commit
	manager.DoTransaction(dto)
	return nil
}
