package GoMybatis

type TransactionStatus struct {
	IsNewTransaction bool
	HasSavepoint     bool
	IsRollbackOnly   bool
	IsCompleted      bool
	Transaction      *Transaction
}

type Transaction struct {
	Id      string
	Sql     string
	Session *Session
}

func (this Transaction) Rollback() error {
	return (*this.Session).Rollback()
}
func (this Transaction) Commit() error {
	return (*this.Session).Commit()
}
func (this Transaction) Begin() error {
	return (*this.Session).Begin()
}

func (this TransactionStatus) Flush() {
	if this.Transaction != nil && this.Transaction.Session != nil {
		(*(*this.Transaction).Session).Close()
		this.Transaction.Session = nil
		this.Transaction = nil
	}
}
