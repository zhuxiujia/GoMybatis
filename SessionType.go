package GoMybatis

type SessionType = int

const (
	SessionType_Default      SessionType = iota
	SessionType_Local
	SessionType_TransationRM
	SessionType_UnKnow
)
