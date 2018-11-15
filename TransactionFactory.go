package GoMybatis

type TransactionFactory struct {
	TransactionStatuss map[string][]TransactionStatus
}

func (this *TransactionFactory) New() *TransactionFactory {
	this.TransactionStatuss = make(map[string][]TransactionStatus)
	return this
}

func (this TransactionFactory) GetTransactionStatus(transactionId string) [] TransactionStatus {
	if transactionId == "" {
		return nil
	}
	var result = this.TransactionStatuss[transactionId]
	if result == nil {
		result = make([]TransactionStatus, 0)
		this.TransactionStatuss[transactionId] = result
	}
	return result
}

func (this *TransactionFactory) SetTransactionStatus(transactionId string, transactions [] TransactionStatus) {
	if transactionId == "" {
		return
	}
	this.TransactionStatuss[transactionId] = transactions
}

func (this *TransactionFactory) Append(transactionId string, transaction TransactionStatus) {
	if transactionId == "" {
		return
	}
	var old = this.GetTransactionStatus(transactionId)
	if old != nil {
		old = append(old, transaction)
		this.SetTransactionStatus(transactionId, old)
	}
}
