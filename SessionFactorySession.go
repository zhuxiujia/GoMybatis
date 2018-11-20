package GoMybatis

import "errors"

type SessionFactorySession struct {
	Session
	SessionHolder *Session
	Factory       *SessionFactory
}

func (this *SessionFactorySession) Id() string {
	if this.SessionHolder == nil{
		return ""
	}
	return (*this.SessionHolder).Id()
}
func (this *SessionFactorySession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	if this.SessionHolder == nil{
		return nil,errors.New("[FactorySession] can not run Id(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Query(sqlorArgs)
}
func (this *SessionFactorySession) Exec(sqlorArgs string) (*Result, error) {
	if this.SessionHolder == nil{
		return nil,errors.New("[FactorySession] can not run Exec(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Exec(sqlorArgs)
}
func (this *SessionFactorySession) Rollback() error {
	if this.SessionHolder == nil{
		return errors.New("[FactorySession] can not run Rollback(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Rollback()
}
func (this *SessionFactorySession) Commit() error {
	if this.SessionHolder == nil{
		return errors.New("[FactorySession] can not run Commit(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Commit()
}
func (this *SessionFactorySession) Begin() error {
	if this.SessionHolder == nil{
		return errors.New("[FactorySession] can not run Begin(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Begin()
}
func (this *SessionFactorySession) Close() {
	var id = this.Id()
	var s = this.Factory.SessionMap[id]
	if s != nil {
		if this.SessionHolder != nil{
			(*this.SessionHolder).Close()
		}
		this.Factory.SessionMap[id] = nil
	}
}
