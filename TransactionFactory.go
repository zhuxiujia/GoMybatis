package GoMybatis

type TransactionFactory struct {
	TransactionStatuss map[string]*TransactionStatus
	SessionFactory     *SessionFactory
}

func (it TransactionFactory) New(SessionFactory *SessionFactory) TransactionFactory {
	it.TransactionStatuss = make(map[string]*TransactionStatus)
	it.SessionFactory = SessionFactory
	return it
}

func (it *TransactionFactory) GetTransactionStatus(mapperName string,transactionId string) (*TransactionStatus, error) {
	var Session Session
	var result = it.TransactionStatuss[transactionId]
	if result == nil {
		Session = it.SessionFactory.NewSession(mapperName,SessionType_Default, nil)
		var transaction = Transaction{
			Id:      transactionId,
			Session: Session,
		}
		var transactionStatus = TransactionStatus{
			IsNewTransaction: true,
			Transaction:      &transaction,
		}
		result = &transactionStatus
		it.TransactionStatuss[transactionId] = result
	}
	return result, nil
}

func (it *TransactionFactory) SetTransactionStatus(transactionId string, transaction *TransactionStatus) {
	if transactionId == "" {
		return
	}
	it.TransactionStatuss[transactionId] = transaction
}

func (it *TransactionFactory) Append(mapperName string,transactionId string, transaction TransactionStatus) {
	if transactionId == "" {
		return
	}
	var old, _ = it.GetTransactionStatus(mapperName,transactionId)
	if old != nil {
		it.SetTransactionStatus(transactionId, old)
	}
}
