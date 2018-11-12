package GoMybatis

import "github.com/zhuxiujia/GoMybatis/lib/github.com/go-xorm/xorm"

type XormEngine struct {
	SqlEngine
	Engine *xorm.Engine
}

func (this XormEngine)Query(sqlorArgs string) ([]map[string][]byte,error){
	return this.Engine.Query(sqlorArgs)
}

func (this XormEngine)Exec(sqlorArgs string) (Result, error){
	var sqlRes,e = this.Engine.Exec(sqlorArgs)
	if e!=nil{
		var LastInsertId,_=sqlRes.LastInsertId()
		var RowsAffected,_=sqlRes.RowsAffected()
		var res = Result{
			LastInsertId:LastInsertId,
			RowsAffected:RowsAffected,
		}
		return res,nil
	}else {
		return Result{},nil
	}
}
//bean 工厂，根据xml配置创建函数,并且动态代理到你定义的struct func里
//bean 参数必须为指针类型,指向你定义的struct
//你定义的struct必须有可导出的func属性,例如：
//type MyUserMapperImpl struct {
//	UserMapper                                                 `mapperPath:"/mapper/user/UserMapper.xml"`
//	SelectById    func(id string, result *model.User) error    `mapperParams:"id"`
//	SelectByPhone func(id string, phone string, result *model.User) error `mapperParams:"id,phone"`
//	DeleteById    func(id string, result *int64) error         `mapperParams:"id"`
//	Insert        func(arg model.User, result *int64) error
//}
//func的参数支持2种函数，第一种函数 基本参数个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//func的参数支持2种函数，第二种函数第一个参数必须为结构体(例如 arg model.User,该结构体的属性可以指定tag `json:"xxx"`为参数名称),最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//使用UseProxyMapper函数设置代理后即可正常使用。
func UseProxyMapperByXorm(bean interface{}, xml []byte, xormEngine *xorm.Engine) {
	var engine = XormEngine{
		Engine:xormEngine,
	}
	var sqlEngine=SqlEngine(engine)
	UseProxyMapper(bean,xml,&sqlEngine)
}