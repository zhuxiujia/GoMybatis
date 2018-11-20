package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"errors"
)

type TransationRMServerSession struct {
	Session
	SessionId string
	Client *TransationRMClient
}

func (this *TransationRMServerSession)Id() string {
	return this.SessionId
}

func (this *TransationRMServerSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	var result TransactionRspDTO
	var error = this.Client.Call(TransactionReqDTO{Status: Transaction_Status_NO, TransactionId: utils.CreateUUID(), Sql: sqlorArgs, ActionType: ActionType_Query,}, &result)
	if error == nil && result.Error != "" {
		error = errors.New(result.Error)
	}
	return result.Query, error
}
func (this *TransationRMServerSession) Exec(sqlorArgs string) (*Result, error) {
	var result TransactionRspDTO
	var error = this.Client.Call(TransactionReqDTO{Status: Transaction_Status_NO, TransactionId: utils.CreateUUID(), Sql: sqlorArgs, ActionType: ActionType_Exec,}, &result)
	if error == nil && result.Error != "" {
		error = errors.New(result.Error)
	}
	return &result.Exec, error
}
func (this *TransationRMServerSession) Rollback() error {
	panic("[RemoteSession] not alow local Rollback()")
	return nil
}
func (this *TransationRMServerSession) Commit() error {
	panic("[RemoteSession] not alow local Commit()")
	return nil
}
func (this *TransationRMServerSession) Begin() error {
	panic("[RemoteSession] not alow local Begin()")
	return nil
}
func (this *TransationRMServerSession) Close() {
}
