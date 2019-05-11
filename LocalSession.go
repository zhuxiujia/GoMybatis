package GoMybatis

import (
	"database/sql"
	"errors"
	"github.com/zhuxiujia/GoMybatis/tx"
	"github.com/zhuxiujia/GoMybatis/utils"
	"strconv"
)

//本地直连session
type LocalSession struct {
	SessionId       string
	driver          string
	url             string
	db              *sql.DB
	stmt            *sql.Stmt
	txStack         tx.TxStack
	savePointStack  *tx.SavePointStack
	isClosed        bool
	newLocalSession *LocalSession

	logSystem Log
}

func (it LocalSession) New(driver string, url string, db *sql.DB, logSystem Log) LocalSession {
	return LocalSession{
		SessionId: utils.CreateUUID(),
		db:        db,
		txStack:   tx.TxStack{}.New(),
		driver:    driver,
		url:       url,
		logSystem: logSystem,
	}
}

func (it *LocalSession) Id() string {
	return it.SessionId
}

func (it *LocalSession) Rollback() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Rollback() a Closed Session!")
	}

	if it.newLocalSession != nil {
		var e = it.newLocalSession.Rollback()
		it.newLocalSession.Close()
		it.newLocalSession = nil
		if e != nil {
			return e
		}
	}

	var t, p = it.txStack.Pop()
	if t != nil && p != nil {
		if *p == tx.PROPAGATION_NESTED {
			if it.savePointStack == nil {
				var stack = tx.SavePointStack{}.New()
				it.savePointStack = &stack
			}
			var point = it.savePointStack.Pop()
			if point != nil {
				if it.logSystem != nil {
					it.logSystem.Println([]byte("[GoMybatis] exec ====================" + "rollback to " + *point))
				}
				_, e := t.Exec("rollback to " + *point)
				e = it.dbErrorPack(e)
				if e != nil {
					return e
				}
			}
		}

		if it.txStack.Len() == 0 {
			if it.logSystem != nil {
				it.logSystem.Println([]byte("[GoMybatis] Rollback tx session:" + it.Id()))
			}
			var err = t.Rollback()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (it *LocalSession) Commit() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Commit() a Closed Session!")
	}

	if it.newLocalSession != nil {
		var e = it.newLocalSession.Commit()
		it.newLocalSession.Close()
		it.newLocalSession = nil
		if e != nil {
			return e
		}
	}

	var t, p = it.txStack.Pop()
	if t != nil && p != nil {

		if *p == tx.PROPAGATION_NESTED {
			if it.savePointStack == nil {
				var stack = tx.SavePointStack{}.New()
				it.savePointStack = &stack
			}
			var pId = "p" + strconv.Itoa(it.txStack.Len()+1)
			it.savePointStack.Push(pId)
			if it.logSystem != nil {
				it.logSystem.Println([]byte("[GoMybatis] exec " + "savepoint " + pId))
			}
			_, e := t.Exec("savepoint " + pId)
			e = it.dbErrorPack(e)
			if e != nil {
				return e
			}
		}
		if it.txStack.Len() == 0 {
			if it.logSystem != nil {
				it.logSystem.Println([]byte("[GoMybatis] Commit tx session:" + it.Id()))
			}
			var err = t.Commit()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (it *LocalSession) Begin(p *tx.Propagation) error {
	var propagation = ""
	if p != nil {
		propagation = tx.ToString(*p)
	}
	if it.logSystem != nil {
		it.logSystem.Println([]byte("[GoMybatis] Begin session:" + it.Id() + ",Propagation:" + propagation))
	}
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Begin() a Closed Session!")
	}

	if p != nil {
		switch *p {
		case tx.PROPAGATION_REQUIRED: //end
			if it.txStack.Len() > 0 {
				it.txStack.Push(it.txStack.Last())
				return nil
			} else {
				var t, err = it.db.Begin()
				err = it.dbErrorPack(err)
				if err == nil {
					it.txStack.Push(t, p)
				}
				return err
			}
			break
		case tx.PROPAGATION_SUPPORTS: //end
			if it.txStack.Len() > 0 {
				return nil
			} else {
				//非事务
				return nil
			}
			break
		case tx.PROPAGATION_MANDATORY: //end
			if it.txStack.Len() > 0 {
				return nil
			} else {
				return errors.New("[GoMybatis] PROPAGATION_MANDATORY Nested transaction exception! current not have a transaction!")
			}
			break
		case tx.PROPAGATION_REQUIRES_NEW:
			if it.txStack.Len() > 0 {
				//TODO stop old tx
			}
			//TODO new session(tx)
			var db, e = sql.Open(it.driver, it.url)
			if e != nil {
				return e
			}
			var sess = LocalSession{}.New(it.driver, it.url, db, it.logSystem) //same PROPAGATION_REQUIRES_NEW
			it.newLocalSession = &sess
			break
		case tx.PROPAGATION_NOT_SUPPORTED:
			if it.txStack.Len() > 0 {
				//TODO stop old tx
			}
			//TODO new session( no tx)
			var db, e = sql.Open(it.driver, it.url)
			if e != nil {
				return e
			}
			var sess = LocalSession{}.New(it.driver, it.url, db, it.logSystem)
			it.newLocalSession = &sess
			break
		case tx.PROPAGATION_NEVER: //END
			if it.txStack.Len() > 0 {
				return errors.New("[GoMybatis] PROPAGATION_NEVER  Nested transaction exception! current Already have a transaction!")
			}
			break
		case tx.PROPAGATION_NESTED: //TODO REQUIRED 类似，增加 save point
			if it.savePointStack == nil {
				var savePointStack = tx.SavePointStack{}.New()
				it.savePointStack = &savePointStack
			}
			if it.txStack.Len() > 0 {
				it.txStack.Push(it.txStack.Last())
				return nil
			} else {
				var tx, err = it.db.Begin()
				err = it.dbErrorPack(err)
				if err == nil {
					it.txStack.Push(tx, p)
				}
				return err
			}
			break
		case tx.PROPAGATION_NOT_REQUIRED: //end
			if it.txStack.Len() > 0 {
				return errors.New("[GoMybatis] PROPAGATION_NOT_REQUIRED Nested transaction exception! current Already have a transaction!")
			} else {
				var tx, err = it.db.Begin()
				err = it.dbErrorPack(err)
				if err == nil {
					it.txStack.Push(tx, p)
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
	if it.logSystem != nil {
		it.logSystem.Println([]byte("[GoMybatis] Close session:" + it.Id()))
	}
	if it.newLocalSession != nil {
		it.newLocalSession.Close()
		it.newLocalSession = nil
	}
	if it.db != nil {
		if it.stmt != nil {
			it.stmt.Close()
		}

		for i := 0; i < it.txStack.Len(); i++ {
			var tx, _ = it.txStack.Pop()
			if tx != nil {
				tx.Rollback()
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
	if it.newLocalSession != nil {
		return it.newLocalSession.Query(sqlorArgs)
	}

	var rows *sql.Rows
	var err error
	var t, _ = it.txStack.Last()
	if t != nil {
		rows, err = t.Query(sqlorArgs)
		err = it.dbErrorPack(err)
	} else {
		rows, err = it.db.Query(sqlorArgs)
		err = it.dbErrorPack(err)
	}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	} else {
		return rows2maps(rows)
	}
	return nil, nil
}

func (it *LocalSession) Exec(sqlorArgs string) (*Result, error) {
	if it.isClosed == true {
		return nil, utils.NewError("LocalSession", " can not Exec() a Closed Session!")
	}
	if it.newLocalSession != nil {
		return it.newLocalSession.Exec(sqlorArgs)
	}

	var result sql.Result
	var err error
	var t, _ = it.txStack.Last()
	if t != nil {
		result, err = t.Exec(sqlorArgs)
		err = it.dbErrorPack(err)
	} else {
		result, err = it.db.Exec(sqlorArgs)
		err = it.dbErrorPack(err)
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

func (it *LocalSession) dbErrorPack(e error) error {
	if e != nil {
		var sqlError = errors.New("[GoMybatis][LocalSession]" + e.Error())
		return sqlError
	}
	return nil
}
