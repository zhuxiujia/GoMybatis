package GoMybatis

import (
	"github.com/zhuxiujia/GoMybatis/utils"
)

type SessionFactorySession struct {
	Session Session
	Factory *SessionFactory
}

func (this *SessionFactorySession) Id() string {
	if this.Session == nil {
		return ""
	}
	return this.Session.Id()
}
func (this *SessionFactorySession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	if this.Session == nil {
		return nil, utils.NewError("SessionFactorySession"," can not run Id(),this.Session == nil")
	}
	return this.Session.Query(sqlorArgs)
}
func (this *SessionFactorySession) Exec(sqlorArgs string) (*Result, error) {
	if this.Session == nil {
		return nil, utils.NewError("SessionFactorySession"," can not run Exec(),this.Session == nil")
	}
	return this.Session.Exec(sqlorArgs)
}
func (this *SessionFactorySession) Rollback() error {
	if this.Session == nil {
		return utils.NewError("SessionFactorySession"," can not run Rollback(),this.Session == nil")
	}
	return this.Session.Rollback()
}
func (this *SessionFactorySession) Commit() error {
	if this.Session == nil {
		return utils.NewError("SessionFactorySession"," can not run Commit(),this.Session == nil")
	}
	return this.Session.Commit()
}
func (this *SessionFactorySession) Begin() error {
	if this.Session == nil {
		return utils.NewError("SessionFactorySession"," can not run Begin(),this.Session == nil")
	}
	return this.Session.Begin()
}
func (this *SessionFactorySession) Close() {
	var id = this.Id()
	var s = this.Factory.SessionMap[id]
	if s != nil {
		if this.Session != nil {
			this.Session.Close()
		}
		this.Factory.SessionMap[id] = nil
	}
}
