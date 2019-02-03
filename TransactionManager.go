package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
	"log"
)

type Transaction_Status int

const (
	Transaction_Status_NO       Transaction_Status = iota //非事务
	Transaction_Status_Prepare                            //准备事务
	Transaction_Status_Commit                             //提交事务
	Transaction_Status_Rollback                           //回滚事务
)

func (status Transaction_Status) ToString() string {
	switch status {
	case Transaction_Status_NO:
		return "Transaction_Status_NO"
	case Transaction_Status_Prepare:
		return "Transaction_Status_Prepare"
	case Transaction_Status_Commit:
		return "Transaction_Status_Commit"
	case Transaction_Status_Rollback:
		return "Transaction_Status_Rollback"
	default:
		return "not init Transaction_Status!"
	}
}

type ActionType int

const (
	ActionType_Exec  ActionType = iota //执行
	ActionType_Query                   //查询
)

type TransactionReqDTO struct {
	Status        Transaction_Status
	TransactionId string //事务id(不可空)
	OwnerId       string //所有者
	Sql           string //sql内容(可空)
	ActionType    ActionType
	MapperName    string //mapper名称
}

type TransactionRspDTO struct {
	TransactionId string //事务id(不可空)
	Error         string
	Success       int
	Query         []map[string][]byte
	Exec          Result
}

type TransactionManager interface {
	GetTransaction(mapperName string, def *TransactionDefinition, transactionId string, OwnerId string) (*TransactionStatus, error)
	Commit(mapperName string, transactionId string) error
	Rollback(mapperName string, transactionId string) error
}

type DefaultTransationManager struct {
	TransactionManager
	SessionFactory     *SessionFactory
	TransactionFactory *TransactionFactory
}

func (it DefaultTransationManager) New(SessionFactory *SessionFactory, TransactionFactory *TransactionFactory) DefaultTransationManager {
	it.SessionFactory = SessionFactory
	it.TransactionFactory = TransactionFactory
	return it
}

func (it DefaultTransationManager) GetTransaction(mapperName string, def *TransactionDefinition, transactionId string, OwnerId string) (*TransactionStatus, error) {
	if transactionId == "" {
		return nil, utils.NewError("TransactionManager", " transactionId ="+transactionId+" transations is nil!")
	}
	if def == nil {
		var d = TransactionDefinition{}.Default()
		def = &d
	}
	//TODO equal mapperName
	var transationStatus, err = it.TransactionFactory.GetTransactionStatus(mapperName, transactionId)
	if err != nil {
		return nil, err
	}
	if def.PropagationBehavior == PROPAGATION_REQUIRED {
		//todo doBegin
		if transationStatus.IsNewTransaction {
			//新事务，则调用begin
			transationStatus.OwnerId = OwnerId
			var err = transationStatus.Begin()
			if err == nil {
				if def.Timeout != 0 {
					//transation out of time,default not set out of time
					//事务超时,时间大于0则启动超时机制
					transationStatus.DelayFlush(def.Timeout)
				}
				return transationStatus, err
			}
		}
	}
	return transationStatus, nil
}

func (it DefaultTransationManager) Commit(mapperName string, transactionId string) error {
	//TODO equal mapperName
	var transactions, err = it.TransactionFactory.GetTransactionStatus(mapperName, transactionId)
	if err != nil {
		log.Println(err)
		return err
	}
	return transactions.Commit()
}

func (it DefaultTransationManager) Rollback(mapperName string, transactionId string) error {
	//TODO equal mapperName
	var transactions, err = it.TransactionFactory.GetTransactionStatus(mapperName, transactionId)
	if err != nil {
		log.Println(err)
		return err
	}
	return transactions.Rollback()
}

//执行事务
func (it DefaultTransationManager) DoTransaction(dto TransactionReqDTO) TransactionRspDTO {
	var transcationStatus *TransactionStatus
	var err error

	transcationStatus, err = it.GetTransaction(dto.MapperName, nil, dto.TransactionId, dto.OwnerId)
	if transcationStatus == nil || transcationStatus.Transaction == nil || transcationStatus.Transaction.Session == nil {
		return TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Error:         "Transaction does not exist,id=" + dto.TransactionId,
		}
	}
	if err != nil {
		return TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Error:         err.Error(),
		}
	}
	if err != nil {
		return TransactionRspDTO{
			TransactionId: dto.TransactionId,
			Error:         err.Error(),
		}
	}
	log.Println("[TransactionManager] do transactionId=", dto.TransactionId, ",sessionId=", transcationStatus.Transaction.Session.Id(), "status=", dto.Status.ToString())

	if dto.Status == Transaction_Status_NO {
		defer transcationStatus.Flush() //关闭
		return it.DoAction(dto, transcationStatus)
	} else if dto.Status == Transaction_Status_Prepare {
		return it.DoAction(dto, transcationStatus)
	} else if dto.Status == Transaction_Status_Commit {
		if transcationStatus.OwnerId == dto.OwnerId { //PROPAGATION_REQUIRED 情况下 子事务 不可提交
			defer transcationStatus.Flush() //关闭
			err = transcationStatus.Commit()
			if err != nil {
				return TransactionRspDTO{
					TransactionId: dto.TransactionId,
					Error:         err.Error(),
				}
			}
			//TODO equal mapperName
			var transaction, err = it.TransactionFactory.GetTransactionStatus(dto.MapperName, dto.TransactionId)
			if err != nil {
				log.Println(err)
			} else {
				transaction.Flush()
			}
		}
	} else if dto.Status == Transaction_Status_Rollback {
		defer transcationStatus.Flush() //关闭，//PROPAGATION_REQUIRED 情况下 子事务 可关闭
		err = transcationStatus.Rollback()
		if err != nil {
			return TransactionRspDTO{
				TransactionId: dto.TransactionId,
				Error:         err.Error(),
			}
		}
	} else {
		err = utils.NewError("TransactionManager", " arg have no action!")
	}
	var errString = ""
	if err != nil {
		errString = err.Error()
	}
	return TransactionRspDTO{
		TransactionId: dto.TransactionId,
		Error:         errString,
	}
}

//执行数据库操作
func (it DefaultTransationManager) DoAction(dto TransactionReqDTO, transcationStatus *TransactionStatus) TransactionRspDTO {
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
		log.Println("[TransactionManager] TransactionId:", dto.TransactionId, ",Exec:", dto.Sql)
		var res, e = transcationStatus.Transaction.Session.Exec(dto.Sql)
		var err string
		if e != nil {
			err = e.Error()
			return TransactionRspDTO{
				TransactionId: dto.TransactionId,
				Error:         err,
			}
		} else {
			return TransactionRspDTO{
				TransactionId: dto.TransactionId,
				Exec:          *res,
				Error:         err,
			}
		}
	} else {
		log.Println("[TransactionManager] Query ", dto.Sql)
		var res, e = transcationStatus.Transaction.Session.Query(dto.Sql)
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
