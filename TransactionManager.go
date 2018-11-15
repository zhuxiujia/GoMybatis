package GoMybatis

import "errors"

type TransactionManager interface {
	GetTransaction(def *TransactionDefinition) (*TransactionStatus, error)
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

func (this DefaultTransationManager) GetTransaction(def *TransactionDefinition, sessionId string, transactionId string) (*TransactionStatus, error) {
	if sessionId == "" {
		return nil, errors.New("[TransactionManager] sessionId can not be nil!")
	}
	if def == nil {
		var d=TransactionDefinition{}.Default()
		def=&d
	}
	var session = this.SessionFactory.GetSession(sessionId)
	if session == nil {
		return nil, errors.New("[TransactionManager] sessionId =" + sessionId + " session is nil!")
	}
	var transation = this.TransactionFactory.GetTransactionStatus(transactionId)
	if transation == nil {
		return nil, errors.New("[TransactionManager] transactionId =" + transactionId + " transations is nil!")
	}
	if def.PropagationBehavior==PROPAGATION_REQUIRED{
        //todo doBegin
		(*transation).Transaction.Begin()
	}

	return transation,nil
}

func (this DefaultTransationManager) Commit() error {
   return nil
}

func (this DefaultTransationManager) Rollback(status TransactionStatus) error {
	return nil
}
