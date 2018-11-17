package GoMybatis

import (
	"database/sql"
	"errors"
)

type Result struct {
	LastInsertId int64
	RowsAffected int64
}

type Session interface {
	Query(sqlorArgs string) ([]map[string][]byte, error)
	Exec(sqlorArgs string) (Result, error)
	Rollback() error
	Commit() error
	Begin() error
	Close()
}

//产生session的引擎
type SessionEngine interface {
	NewSession() *Session
}

//本地直连session
type LocalSqlSession struct {
	Session
	Id                     *string
	db                     *sql.DB
	stmt                   *sql.Stmt
	tx                     *sql.Tx
	isCommitedOrRollbacked *bool
}

func (this *LocalSqlSession) Rollback() error {
	if this.tx != nil {
		var err = this.tx.Rollback()
		if err == nil {
			*this.isCommitedOrRollbacked = true
		} else {
			return err
		}
	}
	return nil
}

func (this *LocalSqlSession) Commit() error {
	if this.tx != nil {
		var err = this.tx.Commit()
		if err == nil {
			*this.isCommitedOrRollbacked = true
		}
	}
	return nil
}

func (this *LocalSqlSession) Begin() error {
	if this.tx == nil {
		var tx, err = this.db.Begin()
		if err == nil {
			this.tx = tx
		} else {
			return err
		}
	}
	return nil
}

func (this *LocalSqlSession) Close() {
	if this.db != nil {
		if this.stmt != nil {
			this.stmt.Close()
		}
		// When Close be called, if session is a transaction and do not call
		// Commit or Rollback, then call Rollback.
		if this.tx != nil && !*this.isCommitedOrRollbacked {
			this.tx.Rollback()
		}
		this.tx = nil
		this.db = nil
		this.stmt = nil
	}
}

func (this *LocalSqlSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	var rows *sql.Rows
	var err error
	if this.tx != nil {
		rows, err = this.tx.Query(sqlorArgs)
	} else {
		rows, err = this.db.Query(sqlorArgs)
	}
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		return rows2maps(rows)
	}
	return nil, nil
}

func (this *LocalSqlSession) Exec(sqlorArgs string) (Result, error) {
	var result sql.Result
	var err error
	if this.tx != nil {
		if *this.isCommitedOrRollbacked {
			return Result{}, errors.New("Exec sql fail!, session isCommitedOrRollbacked!")
		}
		result, err = this.tx.Exec(sqlorArgs)
	} else {
		result, err = this.db.Exec(sqlorArgs)
	}
	if err != nil {
		return Result{}, err
	} else {
		var LastInsertId, _ = result.LastInsertId()
		var RowsAffected, _ = result.RowsAffected()
		return Result{
			LastInsertId: LastInsertId,
			RowsAffected: RowsAffected,
		}, nil
	}
}