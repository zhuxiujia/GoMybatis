package GoMybatis

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
	var factorySession = SessionFactorySession{
		Session: *session,
		Factory: this,
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
