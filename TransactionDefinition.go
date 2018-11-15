package GoMybatis

type PROPAGATION int
type ISOLATION int

const (
	ISOLATION_DEFAULT          ISOLATION = -1
	ISOLATION_READ_UNCOMMITTED ISOLATION = 1
	ISOLATION_READ_COMMITTED   ISOLATION = 2
	ISOLATION_REPEATABLE_READ  ISOLATION = 3
	ISOLATION_SERIALIZABLE     ISOLATION = 4
)

const (
	PROPAGATION_REQUIRED      PROPAGATION = iota
	PROPAGATION_SUPPORTS
	PROPAGATION_MANDATORY
	PROPAGATION_REQUIRES_NEW
	PROPAGATION_NOT_SUPPORTED
	PROPAGATION_NEVER
	PROPAGATION_NESTED
)

type TransactionDefinition struct {
	PropagationBehavior PROPAGATION
	IsolationLevel      ISOLATION
	Timeout             int
	IsReadOnly          bool
}

func (this TransactionDefinition) Default() TransactionDefinition {
	return TransactionDefinition{
		PropagationBehavior: PROPAGATION_REQUIRED,
		IsolationLevel:      ISOLATION_DEFAULT,
		Timeout:             0,
		IsReadOnly:          false,
	}
}
