package GoMybatis

type SessionFactory struct {
	Engine     SessionEngine
	SessionMap map[string]Session
}

func (it SessionFactory) New(Engine SessionEngine) SessionFactory {
	it.Engine = Engine
	it.SessionMap = make(map[string]Session)
	return it
}

func (it *SessionFactory) NewSession(mapperName string, sessionType SessionType) Session {
	if it.SessionMap == nil || it.Engine == nil {
		panic("[GoMybatis] SessionFactory not init! you must call method SessionFactory.New(*)")
	}
	var newSession Session
	var err error
	switch sessionType {
	case SessionType_Default:
		var session, err = it.Engine.NewSession(mapperName)
		if err != nil {
			panic(err)
		}
		var factorySession = SessionFactorySession{
			Session: session,
			Factory: it,
		}
		newSession = Session(&factorySession)
		break
	case SessionType_Local:
		newSession, err = it.Engine.NewSession(mapperName)
		if err != nil {
			panic(err)
		}
		break
	default:
		panic("[GoMybatis] newSession() must have a SessionType!")
	}
	it.SessionMap[newSession.Id()] = newSession
	return newSession
}

func (it *SessionFactory) GetSession(id string) Session {
	return it.SessionMap[id]
}

func (it *SessionFactory) SetSession(id string, session Session) {
	it.SessionMap[id] = session
}

func (it *SessionFactory) Close(id string) {
	if it.SessionMap == nil {
		return
	}
	var s = it.SessionMap[id]
	if s != nil {
		s.Close()
		it.SessionMap[id] = nil
	}
}

func (it *SessionFactory) CloseAll(id string) {
	if it.SessionMap == nil {
		return
	}
	for _, v := range it.SessionMap {
		if v != nil {
			v.Close()
			it.SessionMap[id] = nil
		}
	}
}
