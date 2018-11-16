package GoMybatis

type SessionFactory struct {
	Engine   *SqlEngine
}

func (this SessionFactory) New(Engine *SqlEngine) SessionFactory {
	this.Engine = Engine
	return this
}

func (this SessionFactory) GetSession() *Session {
	return (*this.Engine).NewSession()
}
