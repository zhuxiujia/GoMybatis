package GoMybatis

import (
	"GoMybatis/utils"
	"database/sql"
)

//本地直连session
type LocalSession struct {
	SessionId              string
	db                     *sql.DB
	stmt                   *sql.Stmt
	tx                     *sql.Tx
	isCommitedOrRollbacked bool
	isClosed               bool
}

func (it *LocalSession) Id() string {
	return it.SessionId
}

func (it *LocalSession) Rollback() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Rollback() a Closed Session!")
	}
	if it.tx != nil {
		var err = it.tx.Rollback()
		if err == nil {
			it.isCommitedOrRollbacked = true
		} else {
			return err
		}
	}
	return nil
}

func (it *LocalSession) Commit() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Commit() a Closed Session!")
	}
	if it.tx != nil {
		var err = it.tx.Commit()
		if err == nil {
			it.isCommitedOrRollbacked = true
		}
	}
	return nil
}

func (it *LocalSession) Begin() error {
	if it.isClosed == true {
		return utils.NewError("LocalSession", " can not Begin() a Closed Session!")
	}
	if it.tx == nil {
		var tx, err = it.db.Begin()
		if err == nil {
			it.tx = tx
		} else {
			return err
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
		if it.tx != nil && !it.isCommitedOrRollbacked {
			it.tx.Rollback()
		}
		it.tx = nil
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
	if it.tx != nil {
		rows, err = it.tx.Query(sqlorArgs)
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
	if it.tx != nil {
		if it.isCommitedOrRollbacked {
			return nil, utils.NewError("LocalSession", " Exec() sql fail!, session isCommitedOrRollbacked!")
		}
		result, err = it.tx.Exec(sqlorArgs)
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
