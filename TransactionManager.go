package GoMybatis

import (
	"errors"
	"log"
)

type TransactionDTOStatus = int

const (
	Transaction_Status_Pause    = 0 //暂停
	Transaction_Status_Commit   = 1 //提交事务
	Transaction_Status_Rollback = 2 //回滚事务
)

type ActionType = int

const (
	ActionType_Exec  = 0 //执行
	ActionType_Query = 1 //查询
)

type TransactionReqDTO struct {
	Status        TransactionDTOStatus
	TransactionId string //事务id(不可空)
	Sql           string //sql内容(可空)
	ActionType    ActionType
}

type TransactionRspDTO struct {
	TransactionId string //事务id(不可空)
	Error         string
	Success       int
	Query         []map[string][]byte
	Exec          Result
}

type TransactionManager interface {
	GetTransaction(def *TransactionDefinition, transactionId string) (*TransactionStatus, error)
	Commit() error
	Rollback(status TransactionStatus) error
}

type DefaultTransationManager struct {
	TransactionManager
	SessionFactory     *SessionFactory
	TransactionFactory *TransactionFactory
}

func (this DefaultTransationManager) New(SessionFactory *SessionFactory, TransactionFactory *TransactionFactory) DefaultTransationManager {
	this.SessionFactory = SessionFactory
	this.TransactionFactory = TransactionFactory
	return this
}

func (this DefaultTransationManager) GetTransaction(def *TransactionDefinition, transactionId string, OwnerId string) (*TransactionStatus, error) {
	if transactionId == "" {
		return nil, errors.New("[TransactionManager] transactionId =" + transactionId + " transations is nil!")
	}
	if def == nil {
		var d = TransactionDefinition{}.Default()
		def = &d
	}
	var transationStatus = this.TransactionFactory.GetTransactionStatus(transactionId)
	if def.PropagationBehavior == PROPAGATION_REQUIRED {
		//todo doBegin
		if transationStatus.IsNewTransaction {
			//新事务，则调用begin
			transationStatus.IsNewTransaction = false
			transationStatus.OwnerId = OwnerId
			var err = transationStatus.Begin()
			if err != nil {
				return transationStatus, err
			}
		}
	}
	return transationStatus, nil
}

func (this DefaultTransationManager) Commit() error {

	return nil
}

func (this DefaultTransationManager) Rollback(status TransactionStatus) error {

	return nil
}

//执行事务
func (this DefaultTransationManager) DoTransaction(manager DefaultTransationManager, dto TransactionReqDTO, OwnerId string) TransactionRspDTO {
	if dto.TransactionId == "" {
		return TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Error:         "[TransactionManager] arg TransactionId can no be null!",
		}
	}
	transcationStatus, err := manager.GetTransaction(nil, dto.TransactionId, OwnerId)
	if err != nil {
		return TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Error:         err.Error(),
		}
	}
	if dto.Status == Transaction_Status_Pause {
		return this.DoAction(dto, transcationStatus)
	} else if dto.Status == Transaction_Status_Commit {
		if transcationStatus.OwnerId != OwnerId { //PROPAGATION_REQUIRED 情况下 子事务 不可提交
			err = transcationStatus.Commit()
			manager.TransactionFactory.GetTransactionStatus(dto.TransactionId).Flush()
			if err != nil {
				return TransactionRspDTO{
					TransactionId: dto.TransactionId,
					Error:         err.Error(),
				}
			}
			manager.TransactionFactory.GetTransactionStatus(dto.TransactionId).Flush()
		}
	} else if dto.Status == Transaction_Status_Rollback {
		err = transcationStatus.Rollback()
		manager.TransactionFactory.GetTransactionStatus(dto.TransactionId).Flush()
		if err != nil {
			return TransactionRspDTO{
				TransactionId: dto.TransactionId,
				Error:         err.Error(),
			}
		}
	}
	return TransactionRspDTO{
		TransactionId: dto.TransactionId,
		Error:         "[TransactionManager] arg have no action!",
	}
}

//执行数据库操作
func (this DefaultTransationManager) DoAction(dto TransactionReqDTO, transcationStatus *TransactionStatus) TransactionRspDTO {
	if transcationStatus.IsCompleted {
		var TransactionRspDTO = TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Error:         "[TransactionManager] transaction fail!it is completed!",
		}
		return TransactionRspDTO
	}
	if dto.Sql == "" {
		var TransactionRspDTO = TransactionRspDTO{
			TransactionId: dto.TransactionId,
		}
		return TransactionRspDTO
	}
	if dto.ActionType == ActionType_Exec {
		log.Println("[TransactionManager] Exec ", dto.Sql)
		var res, e = (*transcationStatus.Transaction.Session).Exec(dto.Sql)
		var err string
		if e != nil {
			err = e.Error()
		}
		var TransactionRspDTO = TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Exec:          res,
			Error:         err,
		}
		return TransactionRspDTO
	} else {
		log.Println("[TransactionManager] Query ", dto.Sql)
		var res, e = (*transcationStatus.Transaction.Session).Query(dto.Sql)
		var err string
		if e != nil {
			err = e.Error()
		}
		var TransactionRspDTO = TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Query:         res,
			Error:         err,
		}
		return TransactionRspDTO
	}
}
