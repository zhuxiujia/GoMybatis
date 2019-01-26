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

func (this TransationRMSession) New(TransactionId string, Client *TransationRMClient, Status Transaction_Status) *Session {
	this.OwnerId = utils.CreateUUID()
	this.TransactionId = TransactionId
	this.Client = Client
	this.Status = Status
	var Session = Session(&this)
	return &Session
}

func (this *TransationRMSession) Id() string {
	return this.TransactionId
}

func (this *TransationRMSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	var result TransactionRspDTO
	var error = this.Client.Call(TransactionReqDTO{Status: this.Status, TransactionId: this.TransactionId, Sql: sqlorArgs, ActionType: ActionType_Query, OwnerId: this.OwnerId}, &result)
	if error == nil && result.Error != "" {
		error = utils.NewError("TransationRMSession", result.Error)
	}
	return result.Query, error
}

func (this *TransationRMSession) Exec(sqlorArgs string) (*Result, error) {
	var result TransactionRspDTO
	var error = this.Client.Call(TransactionReqDTO{Status: this.Status, TransactionId: this.TransactionId, Sql: sqlorArgs, ActionType: ActionType_Exec, OwnerId: this.OwnerId}, &result)
	if error == nil && result.Error != "" {
		error = utils.NewError("TransationRMSession", result.Error)
	}
	return &result.Exec, error
}

func (this *TransationRMSession) Rollback() error {
	this.Status = Transaction_Status_Rollback
	var result TransactionRspDTO
	return this.Client.Call(TransactionReqDTO{Status: this.Status, TransactionId: this.TransactionId, ActionType: ActionType_Exec, OwnerId: this.OwnerId}, &result)
}

func (this *TransationRMSession) Commit() error {
	this.Status = Transaction_Status_Commit
	var result TransactionRspDTO
	return this.Client.Call(TransactionReqDTO{Status: this.Status, TransactionId: this.TransactionId, ActionType: ActionType_Exec, OwnerId: this.OwnerId}, &result)
}

func (this *TransationRMSession) Begin() error {
	this.Status = Transaction_Status_Prepare
	var result TransactionRspDTO
	var err = this.Client.Call(TransactionReqDTO{Status: this.Status, TransactionId: this.TransactionId, ActionType: ActionType_Exec, OwnerId: this.OwnerId}, &result)
	return err
}

func (this *TransationRMSession) Close() {
	if this.Status == Transaction_Status_Prepare {
		this.Rollback()
	}
	if this.Client != nil {
		this.Client.Close()
	}
}
