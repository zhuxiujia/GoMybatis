package GoMybatis

type TransactionStatus struct {
	OwnerId  string
	IsNewTransaction bool
	HasSavepoint     bool
	IsRollbackOnly   bool
	IsCompleted      bool
	Transaction      *Transaction
}

type Transaction struct {
	Id      string
	Sqls    []string
	Session *Session
}


func (this *TransactionStatus) Append(sql string) error {
	 this.Transaction.Sqls=append(this.Transaction.Sqls, sql)
	 return nil
}

func (this *TransactionStatus) Rollback() error {
	this.IsCompleted=true
	return (*this.Transaction.Session).Rollback()
}

func (this *TransactionStatus) Commit() error {
	this.IsCompleted=true
	return (*this.Transaction.Session).Commit()
}

func (this *TransactionStatus) Begin() error {
	this.IsNewTransaction=false
	return (*this.Transaction.Session).Begin()
}

func (this *TransactionStatus) Flush() {
	if this.Transaction != nil && this.Transaction.Session != nil {
		(*(*this.Transaction).Session).Close()
		this.Transaction.Session = nil
		this.Transaction = nil
	}
}
