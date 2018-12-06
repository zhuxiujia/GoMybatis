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

func (this *TransactionStatus) Rollback() error {
	if this.IsCompleted == true {
		return errors.New("[TransactionManager] can not Rollback() a completed Transaction!")
	}
	this.IsCompleted = true
	defer this.Flush() //close session
	return this.Transaction.Session.Rollback()
}

func (this *TransactionStatus) Commit() error {
	if this.IsCompleted == true {
		return errors.New("[TransactionManager] can not Commit() a completed Transaction!")
	}
	this.IsCompleted = true
	defer this.Flush() //close session
	return this.Transaction.Session.Commit()
}

func (this *TransactionStatus) Begin() error {
	if this.IsNewTransaction == false {
		return errors.New("[TransactionManager] can not Begin() a old Transaction!")
	}
	this.IsNewTransaction = false
	return this.Transaction.Session.Begin()
}

func (this *TransactionStatus) Flush() {
	if this.Transaction != nil && this.Transaction.Session != nil {
		this.Transaction.Session.Close()
		this.Transaction.Session = nil
		this.Transaction = nil
	}
}

//延迟关闭
func (this *TransactionStatus) DelayFlush(t time.Duration) {
	if this.HasSetDelayClose == false {
		go func() {
			time.Sleep(t)
			if this.IsCompleted == false {
				this.Rollback()
			}
		}()
		this.HasSetDelayClose = true
	}
}
