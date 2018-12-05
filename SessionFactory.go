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

func (this *SessionFactory) NewSession(sessionType SessionType, config *TransationRMClientConfig) *Session {
	if this.SessionMap == nil || this.Engine == nil {
		panic("[GoMybatis] SessionFactory not init! you must call method SessionFactory.New(*)")
	}
	var newSession Session
	switch sessionType {
	case SessionType_Default:
		var session = (*this.Engine).NewSession()
		var factorySession = SessionFactorySession{
			Session: *session,
			Factory: this,
		}
		newSession = Session(&factorySession)
		break
	case SessionType_Local:
		var session = (*this.Engine).NewSession()
		newSession = *session
		break
	case SessionType_TransationRM:
		if config == nil {
			panic("[GoMybatis] SessionFactory can not create TransationRMSession,config *TransationRMClientConfig is nil!")
		}
		var transationRMSession = TransationRMSession{}.New(config.TransactionId, &TransationRMClient{
			RetryTime: config.RetryTime,
			Addr:      config.Addr,
		}, config.Status)
		newSession = Session(*transationRMSession)
		break
	default:
		panic("[GoMybatis] newSession() must have a SessionType!")
	}
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
