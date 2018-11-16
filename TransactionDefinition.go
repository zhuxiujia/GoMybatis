package GoMybatis

import "time"

type PROPAGATION int
type ISOLATION int

//隔离级别
const (
	ISOLATION_DEFAULT          ISOLATION = iota - 1
	ISOLATION_READ_UNCOMMITTED
	ISOLATION_READ_COMMITTED
	ISOLATION_REPEATABLE_READ
	ISOLATION_SERIALIZABLE
)

//传播行为
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
	Timeout             time.Duration
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
