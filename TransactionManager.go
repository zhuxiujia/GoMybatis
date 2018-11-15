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

func (this *DefaultTransationManager) New(SessionFactory *SessionFactory, TransactionFactory *TransactionFactory) *DefaultTransationManager {
	this.SessionFactory = SessionFactory
	this.TransactionFactory = TransactionFactory
	return this
}

func (this DefaultTransationManager) GetTransaction(def *TransactionDefinition, sessionId string, transactionId string) (*TransactionStatus, error) {
	if sessionId == "" {
		return nil, errors.New("[TransactionManager] sessionId can not be nil!")
	}
	if def == nil {
		def = TransactionDefinition{}.Default()
	}
	var session = this.SessionFactory.GetSession(sessionId)
	if session == nil {
		return nil, errors.New("[TransactionManager] sessionId =" + sessionId + " session is nil!")
	}
	var transations = this.TransactionFactory.GetTransactionStatus(transactionId)
	if transations == nil {
		return nil, errors.New("[TransactionManager] transactionId =" + transactionId + " transations is nil!")
	}
	if def.PropagationBehavior==PROPAGATION_REQUIRED{
        //todo doBegin

	}

}

func (this DefaultTransationManager) Commit() error {

}

func (this DefaultTransationManager) Rollback(status TransactionStatus) error {

}
