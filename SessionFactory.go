package GoMybatis

type SessionFactory struct {
	Engine   *SessionEngine
}

func (this SessionFactory) New(Engine *SessionEngine) SessionFactory {
	this.Engine = Engine
	return this
}

func (this SessionFactory) GetSession() *Session {
	return (*this.Engine).NewSession()
}
