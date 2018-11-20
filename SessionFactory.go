package GoMybatis

import "errors"

type FactorySession struct {
	Session
	SessionHolder *Session
	Factory       *SessionFactory
}

func (this *FactorySession) Id() string {
	if this.SessionHolder == nil{
		return ""
	}
	return (*this.SessionHolder).Id()
}
func (this *FactorySession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	if this.SessionHolder == nil{
		return nil,errors.New("[FactorySession] can not run Id(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Query(sqlorArgs)
}
func (this *FactorySession) Exec(sqlorArgs string) (*Result, error) {
	if this.SessionHolder == nil{
		return nil,errors.New("[FactorySession] can not run Exec(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Exec(sqlorArgs)
}
func (this *FactorySession) Rollback() error {
	if this.SessionHolder == nil{
		return errors.New("[FactorySession] can not run Rollback(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Rollback()
}
func (this *FactorySession) Commit() error {
	if this.SessionHolder == nil{
		return errors.New("[FactorySession] can not run Commit(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Commit()
}
func (this *FactorySession) Begin() error {
	if this.SessionHolder == nil{
		return errors.New("[FactorySession] can not run Begin(),this.SessionHolder == nil")
	}
	return (*this.SessionHolder).Begin()
}
func (this *FactorySession) Close() {
	var id = this.Id()
	var s = this.Factory.SessionMap[id]
	if s != nil {
		if this.SessionHolder != nil{
			(*this.SessionHolder).Close()
		}
		this.Factory.SessionMap[id] = nil
	}
}

type SessionFactory struct {
	Engine     *SessionEngine
	SessionMap map[string]*Session
}

func (this SessionFactory) New(Engine *SessionEngine) SessionFactory {
	this.Engine = Engine
	this.SessionMap = make(map[string]*Session)
	return this
}

func (this *SessionFactory) NewSession() *Session {
	var session = (*this.Engine).NewSession()
	var factorySession = FactorySession{
		SessionHolder: session,
		Factory:       this,
	}
	var newSession = Session(&factorySession)
	this.SessionMap[newSession.Id()] = &newSession
	return &newSession
}

func (this *SessionFactory) GetSession(id string) *Session {
	return this.SessionMap[id]
}

func (this *SessionFactory) CloseSession(id string) {
	var s = this.SessionMap[id]
	if s != nil {
		(*s).Close()
		this.SessionMap[id] = nil
	}
}
