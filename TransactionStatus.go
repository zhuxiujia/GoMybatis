package GoMybatis

import (
	"errors"
	"time"
)

type TransactionStatus struct {
	OwnerId          string       //所有者id
	IsNewTransaction bool         //是否新启动的事务
	HasSavepoint     bool         //是否保存点
	IsRollbackOnly   bool         //是否只允许rollback
	IsCompleted      bool         //是否完成
	HasSetDelayClose bool         //是否设置了延迟关闭/回滚
	Transaction      *Transaction //事务对象
}

type Transaction struct {
	Id      string
	Session Session
}

func (it *TransactionStatus) Rollback() error {
	if it.IsCompleted == true {
		return errors.New("[TransactionManager] can not Rollback() a completed Transaction!")
	}
	it.IsCompleted = true
	defer it.Flush() //close session
	return it.Transaction.Session.Rollback()
}

func (it *TransactionStatus) Commit() error {
	if it.IsCompleted == true {
		return errors.New("[TransactionManager] can not Commit() a completed Transaction!")
	}
	it.IsCompleted = true
	defer it.Flush() //close session
	return it.Transaction.Session.Commit()
}

func (it *TransactionStatus) Begin() error {
	if it.IsNewTransaction == false {
		return errors.New("[TransactionManager] can not Begin() a old Transaction!")
	}
	it.IsNewTransaction = false
	return it.Transaction.Session.Begin()
}

func (it *TransactionStatus) Flush() {
	if it.Transaction != nil && it.Transaction.Session != nil {
		it.Transaction.Session.Close()
		it.Transaction.Session = nil
		it.Transaction = nil
	}
}

//延迟关闭
func (it *TransactionStatus) DelayFlush(t time.Duration) {
	if it.HasSetDelayClose == false {
		go func() {
			time.Sleep(t)
			if it.IsCompleted == false {
				it.Rollback()
			}
		}()
		it.HasSetDelayClose = true
	}
}
