package GoMybatis

type SessionFactory struct {
	Engine     *SessionEngine
	SessionMap map[string]*Session
}

func (it SessionFactory) New(Engine *SessionEngine) SessionFactory {
	it.Engine = Engine
	it.SessionMap = make(map[string]*Session)
	return it
}

func (it *SessionFactory) NewSession(sessionType SessionType, config *TransationRMClientConfig) Session {
	if it.SessionMap == nil || it.Engine == nil {
		panic("[GoMybatis] SessionFactory not init! you must call method SessionFactory.New(*)")
	}
	var newSession Session
	switch sessionType {
	case SessionType_Default:
		var session = (*it.Engine).NewSession()
		var factorySession = SessionFactorySession{
			Session: session,
			Factory: it,
		}
		newSession = Session(&factorySession)
		break
	case SessionType_Local:
		newSession = (*it.Engine).NewSession()
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
	it.SessionMap[newSession.Id()] = &newSession
	return newSession
}

func (it *SessionFactory) GetSession(id string) *Session {
	return it.SessionMap[id]
}

func (it *SessionFactory) SetSession(id string, session *Session) {
	it.SessionMap[id] = session
}

func (it *SessionFactory) Close(id string) {
	if it.SessionMap == nil {
		return
	}
	var s = it.SessionMap[id]
	if s != nil {
		(*s).Close()
		it.SessionMap[id] = nil
	}
}

func (it *SessionFactory) CloseAll(id string) {
	if it.SessionMap == nil {
		return
	}
	for _, v := range it.SessionMap {
		if v != nil {
			(*v).Close()
			it.SessionMap[id] = nil
		}
	}
}
