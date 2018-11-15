package GoMybatis

type SessionFactory struct {
	Sessions map[string]*Session
	Engine   *SqlEngine
}

func (this SessionFactory) New(Engine *SqlEngine) SessionFactory {
	this.Sessions=make(map[string]*Session)
	this.Engine = Engine
	return this
}

func (this SessionFactory) GetSession(id string) *Session {
	if id == "" {
		return nil
	}
	var session = this.Sessions[id]
	if session == nil {
		session = (*this.Engine).NewSession()
		this.Sessions[id] = session
	}
	return session
}
