package GoMybatis

type TransactionFactory struct {
	TransactionStatuss map[string]*TransactionStatus
	SessionFactory     *SessionFactory
}

func (this TransactionFactory) New(SessionFactory *SessionFactory) TransactionFactory {
	this.TransactionStatuss = make(map[string]*TransactionStatus)
	this.SessionFactory = SessionFactory
	return this
}

func (this *TransactionFactory) GetTransactionStatus(transactionId string) *TransactionStatus {
	var Session *Session
	if transactionId == "" {
		Session = this.SessionFactory.NewSession()
		transactionId = (*Session).Id()
	}
	var result = this.TransactionStatuss[transactionId]
	if result == nil {
		var transaction = Transaction{
			Id:      transactionId,
			Session: Session,
		}
		var transactionStatus = TransactionStatus{
			IsNewTransaction: true,
			Transaction:      &transaction,
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
