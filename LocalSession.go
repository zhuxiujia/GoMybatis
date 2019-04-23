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
	SessionId      string
	driver         string
	url            string
	db             *sql.DB
	stmt           *sql.Stmt
	txStack        tx.TxStack
	savePointStack *tx.SavePointStack
	isClosed       bool

	newLocalSession *LocalSession
}

func (it LocalSession) New(driver string, url string, db *sql.DB) LocalSession {
	return LocalSession{
		SessionId: utils.CreateUUID(),
		db:        db,
		txStack:   tx.TxStack{}.New(),
		driver:    driver,
		url:       url,
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
			var point = it.savePointStack.Pop()
			if point != nil {
				println("[GoMybatis] exec ====================" + "rollback to " + *point)
				r, e := t.Exec("rollback to " + *point)
				println(r)
				if e != nil {
					return e
				}
			}
		}

		if it.txStack.Len() == 0 {
			println("Rollback tx session:", it.Id())
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
			var pId = "p" + strconv.Itoa(it.txStack.Len()+1)
			it.savePointStack.Push(pId)
			println("[GoMybatis]==================== exec " + "savepoint " + pId)
			_, e := t.Exec("savepoint " + pId)
			if e != nil {
				return e
			}
		}
		if it.txStack.Len() == 0 {
			println("Commit tx session:", it.Id())
			var err = t.Commit()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (it *LocalSession) Begin(p *tx.Propagation) error {
	var prog = ""
	if p != nil {
		prog = tx.ToString(*p)
	}
	println("Begin session:", it.Id(), ",prog:", prog)
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
			var sess = LocalSession{}.New(it.driver, it.url, db) //same PROPAGATION_REQUIRES_NEW
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
			var sess = LocalSession{}.New(it.driver, it.url, db)
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
	println("Close session:", it.Id())
	if it.newLocalSession != nil {
		it.newLocalSession.Close()
		it.newLocalSession = nil
	}
	if it.db != nil {
		if it.stmt != nil {
			it.stmt.Close()
		}

		for {
			var tx, _ = it.txStack.Pop()
			if tx != nil {
				tx.Rollback()
			} else {
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
	if it.newLocalSession != nil {
		return it.newLocalSession.Query(sqlorArgs)
	}

	var rows *sql.Rows
	var err error
	var t, _ = it.txStack.Last()
	if t != nil {
		rows, err = t.Query(sqlorArgs)
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
	if it.newLocalSession != nil {
		return it.newLocalSession.Exec(sqlorArgs)
	}

	var result sql.Result
	var err error
	var t, _ = it.txStack.Last()
	if t != nil {
		result, err = t.Exec(sqlorArgs)
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
