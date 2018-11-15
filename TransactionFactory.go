package GoMybatis

type TransactionFactory struct {
	TransactionStatuss map[string]*TransactionStatus
}

func (this TransactionFactory) New() TransactionFactory {
	this.TransactionStatuss = make(map[string]*TransactionStatus)
	return this
}

func (this TransactionFactory) GetTransactionStatus(transactionId string) *TransactionStatus {
	if transactionId == "" {
		return nil
	}
	var result = this.TransactionStatuss[transactionId]
	if result == nil {
		var transaction = Transaction{
			Sqls: make([]string, 0),
		}
		var transactionStatus = TransactionStatus{
			Transaction: &transaction,
		}
		result = &transactionStatus
		this.TransactionStatuss[transactionId] = result
	}
	return result
}

func (this *TransactionFactory) SetTransactionStatus(transactionId string, transaction *TransactionStatus) {
	if transactionId == "" {
		return
	}
	this.TransactionStatuss[transactionId] = transaction
}

func (this *TransactionFactory) Append(transactionId string, transaction TransactionStatus) {
	if transactionId == "" {
		return
	}
	var old = this.GetTransactionStatus(transactionId)
	if old != nil {
		this.SetTransactionStatus(transactionId, old)
	}
}
