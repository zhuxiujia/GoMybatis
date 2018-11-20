package GoMybatis

type FactorySession struct {
	Session
	SessionHolder *Session
	Factory       *SessionFactory
}

func (this *FactorySession) Id() string {
	return (*this.SessionHolder).Id()
}
func (this *FactorySession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	return (*this.SessionHolder).Query(sqlorArgs)
}
func (this *FactorySession) Exec(sqlorArgs string) (Result, error) {
	return (*this.SessionHolder).Exec(sqlorArgs)
}
func (this *FactorySession) Rollback() error {
	return (*this.SessionHolder).Rollback()
}
func (this *FactorySession) Commit() error {
	return (*this.SessionHolder).Commit()
}
func (this *FactorySession) Begin() error {
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
