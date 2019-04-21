package tx

//事务传播行为
type Propagation int

const (
	PROPAGATION_REQUIRED      = iota //默认，表示如果当前事务存在，则支持当前事务。否则，会启动一个新的事务。xorm中默认事务类型。
	PROPAGATION_SUPPORTS             //表示如果当前事务存在，则支持当前事务，如果当前没有事务，就以非事务方式执行。
	PROPAGATION_MANDATORY            //表示如果当前事务存在，则支持当前事务，如果当前没有事务，则返回事务嵌套错误。
	PROPAGATION_REQUIRES_NEW         //表示新建一个全新Session开启一个全新事务，如果当前存在事务，则把当前事务挂起。
	PROPAGATION_NOT_SUPPORTED        //表示以非事务方式执行操作，如果当前存在事务，则新建一个Session以非事务方式执行操作，把当前事务挂起。
	PROPAGATION_NEVER                //表示以非事务方式执行操作，如果当前存在事务，则返回事务嵌套错误。
	PROPAGATION_NESTED               //表示如果当前事务存在，则在嵌套事务内执行，如嵌套事务回滚，则只会在嵌套事务内回滚，不会影响当前事务。如果当前没有事务，则进行与PROPAGATION_REQUIRED类似的操作。
	PROPAGATION_NOT_REQUIRED         //表示如果当前没有事务，就新建一个事务,否则返回错误。
)
