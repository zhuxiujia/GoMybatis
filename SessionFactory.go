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
	this.SessionMap[(*session).Id()] = session
	return session
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
