package GoMybatis

type SessionType = int

const (
	SessionType_Default      SessionType = iota //默认session类型
	SessionType_Local                           //本地session
	SessionType_TransationRM                    //远程session
)
