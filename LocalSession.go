package GoMybatis

import (
	"database/sql"
	"errors"
	"github.com/zhuxiujia/GoMybatis/tx"
	"github.com/zhuxiujia/GoMybatis/utils"
)

//本地直连session
type LocalSession struct {
	SessionId   string
	db          *sql.DB
	stmt        *sql.Stmt
	txStack     tx.TxStack
	isClosed    bool
	propagation *tx.Propagation

	newLocalSession *LocalSession
}

func (it LocalSession) New(db *sql.DB, propagation *tx.Propagation) LocalSession {
	return LocalSession{
		SessionId:   utils.CreateUUID(),
		db:          db,
		txStack:     tx.TxStack{}.New(),
		propagation: propagation,
	}
}

func (it *LocalSession) Id() string {
	return it.SessionId
}

func (it *LocalSession) Rollback() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Rollback() a Closed Session!")
	}
	var tx = it.txStack.Pop()
	if tx != nil {
		var err = tx.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}

func (it *LocalSession) Commit() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Commit() a Closed Session!")
	}
	var tx = it.txStack.Pop()
	if tx != nil {
		var err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (it *LocalSession) Begin() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Begin() a Closed Session!")
	}

	if it.propagation != nil {
		switch *it.propagation {
		case tx.PROPAGATION_REQUIRED://end
			if it.txStack.Len() > 0 {
				return nil
			} else {
				var tx, err = it.db.Begin()
				if err == nil {
					it.txStack.Push(tx)
				}
				return err
			}
			break
		case tx.PROPAGATION_SUPPORTS://end
			if it.txStack.Len() > 0 {
				return nil
			} else {
				//非事务
				return nil
			}
			break
		case tx.PROPAGATION_MANDATORY://end
			if it.txStack.Len() > 0 {
				return nil
			} else {
				return errors.New("[GoMybatis] PROPAGATION_MANDATORY Nested transaction exception! current not have a transaction!")
			}
			break
		case tx.PROPAGATION_REQUIRES_NEW://TODO
			if it.txStack.Len() > 0 {
				//TODO stop old tx
			}
			//TODO new session(tx)
			break
		case tx.PROPAGATION_NOT_SUPPORTED://TODO
			if it.txStack.Len() > 0 {
				//TODO stop old tx
			}
			//TODO new session( no tx)
			break
		case tx.PROPAGATION_NEVER://END
			if it.txStack.Len() > 0 {
				return errors.New("[GoMybatis] PROPAGATION_NEVER  Nested transaction exception! current Already have a transaction!")
			}
			break
		case tx.PROPAGATION_NESTED: //TODO REQUIRED 类似，增加 save point
			if it.txStack.Len() > 0 {
				return nil
			} else {
				var tx, err = it.db.Begin()
				if err == nil {
					it.txStack.Push(tx)
				}
				return err
			}
			break
		case tx.PROPAGATION_NOT_REQUIRED://end
			if it.txStack.Len() > 0 {
				return errors.New("[GoMybatis] PROPAGATION_NOT_REQUIRED Nested transaction exception! current Already have a transaction!")
			} else {
				var tx, err = it.db.Begin()
				if err == nil {
					it.txStack.Push(tx)
				}
				return err
			}
			break
		default:
			panic("[GoMybatis] Nested transaction exception! not support PROPAGATION in begin!")
			break
		}

	}
	return nil
}

func (it *LocalSession) Close() {
	if it.db != nil {
		if it.stmt != nil {
			it.stmt.Close()
		}
		// When Close be called, if session is a transaction and do not call
		// Commit or Rollback, then call Rollback.

		for {
			var tx = it.txStack.Pop()
			tx.Rollback()
			if tx == nil {
				break
			}
		}

		it.db = nil
		it.stmt = nil
		it.isClosed = true
	}
}

func (it *LocalSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	if it.isClosed == true {
		return nil, utils.NewError("LocalSession", " can not Query() a Closed Session!")
	}
	var rows *sql.Rows
	var err error
	if it.txStack.Last() != nil {
		rows, err = it.txStack.Last().Query(sqlorArgs)
	} else {
		rows, err = it.db.Query(sqlorArgs)
	}
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		return rows2maps(rows)
	}
	return nil, nil
}

func (it *LocalSession) Exec(sqlorArgs string) (*Result, error) {
	if it.isClosed == true {
		return nil, utils.NewError("LocalSession", " can not Exec() a Closed Session!")
	}
	var result sql.Result
	var err error
	if it.txStack.Last() != nil {
		result, err = it.txStack.Last().Exec(sqlorArgs)
	} else {
		result, err = it.db.Exec(sqlorArgs)
	}
	if err != nil {
		return nil, err
	} else {
		var LastInsertId, _ = result.LastInsertId()
		var RowsAffected, _ = result.RowsAffected()
		return &Result{
			LastInsertId: LastInsertId,
			RowsAffected: RowsAffected,
		}, nil
	}
}
