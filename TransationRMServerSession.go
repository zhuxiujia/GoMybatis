package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
)

type TransationRMSession struct {
	TransactionId string
	OwnerId       string
	Client        *TransationRMClient
	Status        Transaction_Status //默认0，非事务
}

func (it TransationRMSession) New(TransactionId string, Client *TransationRMClient, Status Transaction_Status) *Session {
	it.OwnerId = utils.CreateUUID()
	it.TransactionId = TransactionId
	it.Client = Client
	it.Status = Status
	var Session = Session(&it)
	return &Session
}

func (it *TransationRMSession) Id() string {
	return it.TransactionId
}

func (it *TransationRMSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	var result TransactionRspDTO
	var error = it.Client.Call(TransactionReqDTO{Status: it.Status, TransactionId: it.TransactionId, Sql: sqlorArgs, ActionType: ActionType_Query, OwnerId: it.OwnerId}, &result)
	if error == nil && result.Error != "" {
		error = utils.NewError("TransationRMSession", result.Error)
	}
	return result.Query, error
}

func (it *TransationRMSession) Exec(sqlorArgs string) (*Result, error) {
	var result TransactionRspDTO
	var error = it.Client.Call(TransactionReqDTO{Status: it.Status, TransactionId: it.TransactionId, Sql: sqlorArgs, ActionType: ActionType_Exec, OwnerId: it.OwnerId}, &result)
	if error == nil && result.Error != "" {
		error = utils.NewError("TransationRMSession", result.Error)
	}
	return &result.Exec, error
}

func (it *TransationRMSession) Rollback() error {
	it.Status = Transaction_Status_Rollback
	var result TransactionRspDTO
	return it.Client.Call(TransactionReqDTO{Status: it.Status, TransactionId: it.TransactionId, ActionType: ActionType_Exec, OwnerId: it.OwnerId}, &result)
}

func (it *TransationRMSession) Commit() error {
	it.Status = Transaction_Status_Commit
	var result TransactionRspDTO
	return it.Client.Call(TransactionReqDTO{Status: it.Status, TransactionId: it.TransactionId, ActionType: ActionType_Exec, OwnerId: it.OwnerId}, &result)
}

func (it *TransationRMSession) Begin() error {
	it.Status = Transaction_Status_Prepare
	var result TransactionRspDTO
	var err = it.Client.Call(TransactionReqDTO{Status: it.Status, TransactionId: it.TransactionId, ActionType: ActionType_Exec, OwnerId: it.OwnerId}, &result)
	return err
}

func (it *TransationRMSession) Close() {
	if it.Status == Transaction_Status_Prepare {
		it.Rollback()
	}
	if it.Client != nil {
		it.Client.Close()
	}
}
