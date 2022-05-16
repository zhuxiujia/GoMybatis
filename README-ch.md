# 更简单的xml/sql风格 功能丰富 值得一试的ORM库

## [网站](https://zhuxiujia.github.io/gomybatis.io/#/)


[![Go Report Card](https://goreportcard.com/badge/github.com/zhuxiujia/GoMybatis)](https://goreportcard.com/report/github.com/zhuxiujia/GoMybatis)
[![Build Status](https://travis-ci.com/zhuxiujia/GoMybatis.svg?branch=master)](https://travis-ci.com/zhuxiujia/GoMybatis)
[![GoDoc](https://godoc.org/github.com/zhuxiujia/GoMybatis?status.svg)](https://godoc.org/github.com/zhuxiujia/GoMybatis)
[![Coverage Status](https://coveralls.io/repos/github/zhuxiujia/GoMybatis/badge.svg?branch=master)](https://coveralls.io/github/zhuxiujia/GoMybatis?branch=master)
[![codecov](https://codecov.io/gh/zhuxiujia/GoMybatis/branch/master/graph/badge.svg)](https://codecov.io/gh/zhuxiujia/GoMybatis)


![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/vuetify.png)

### 实际应用[点击-GoMybatis商城](https://github.com/zhuxiujia/GoMybatisMall)

### 请仔细阅读[网站](https://zhuxiujia.github.io/gomybatis.io/#/)



# future
* <a href="https://zhuxiujia.github.io/gomybatis.io/">稳定</a>，已应用生产环境App（电商/金融/卡充值类），功能稳定，适合各类 大小型项目以及复杂的金融项目,ERP项目 帮助您将数十万RMB轻松收入囊中<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">高性能</a>，单机最高可达751020 Qps/s,总耗时0.14s （测试环境返回模拟sql数据，并发1000，总数100000，6核16GB win10）<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">无痛迁移</a>，xml对于语言无关/低耦合，兼容大部分Java(Mybatis3,Mybatis Plus)框架逻辑，无痛苦Java Spring Mybatis的xml sql文件迁移至Go语言（仅修改resultMap的javaType为langType指定go语言类型）<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">声明式事务、AOP事务、事务传播行为</a>只需在函数尾部 定义`tx:"" rollback:"error"`即可启用声明式事务，事务传播行为,回滚策略.轻松应对复杂的事务嵌套/回滚<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">扩展日志接口</a>异步消息队列日,框架内sql日志使用带缓存的channel实现 消息队列异步记录日志<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">动态SQL</a>，在xml中`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<where>,<foreach>,<resultMap>,<bind>,<choose><when><otherwise>,<sql><include>`等等java框架Mybatis包含的15种实用功能<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">智能表达式</a>可处理动态判断、计算任务（例如：`#{foo.Bar}#{arg+1}#{arg*1}#{arg/1}#{arg-1}`）,例如写模糊查询`select * from table where phone like #{phone+'%'}`(注意后置百分号走索引)<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">动态数据源</a>可自定义多数据源，动态切换多个数据库实例<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">模板标签（新）</a>一行代码实现增删改查，逻辑删除，乐观锁，而且还保留完美的扩展性（标签体内可以继续扩展sql逻辑）<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">乐观锁（新）</a>`<updateTemplate>`乐观锁,尽可能防止并发竞争修改记录<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">逻辑删除（新）</a>`<insertTemplate><updateTemplate><deleteTemplate><selectTemplate>`逻辑删除,防止意外删除数据，数据恢复简单<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">RPC/MVC组件支持（新）</a>让服务完美支持RPC（减少参数限制）,动态代理，事务订阅，易于微服务集成和扩展 详情请点击链接https://github.com/zhuxiujia/easyrpc<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Wrapper插件免写sql支持，Page分页插件支持（预计在golang1.17版本（支持泛型）之后出现）<br>

## 数据库驱动(支持所有基于标准库database/sql下的所有驱动)
``` bash
 //传统数据库
 Mysql:                             github.com/go-sql-driver/mysql
 MyMysql:                           github.com/ziutek/mymysql/godrv
 Postgres:                          github.com/lib/pq
 SQLite:                            github.com/mattn/go-sqlite3
 MsSql:                             github.com/denisenkom/go-mssqldb
 Oracle:                            github.com/mattn/go-oci8
 //分布式NewSql数据库
 Tidb:                              github.com/go-sql-driver/mysql
 CockroachDB:                       github.com/lib/pq
 ```
 
## 使用教程
> 教程源码  https://github.com/zhuxiujia/GoMybatis/tree/master/example

* GoPath使用： go get 命令下载GoMybatis和对应的数据库驱动
``` bash
go get github.com/zhuxiujia/GoMybatis
```
``` bash
//驱动
go get github.com/go-sql-driver/mysql
```
* mod使用（环境变量加入GO111MODULE auto）:
``` bash
//命令行执行 禁止gosumdb
go env -w GOSUMDB=off
```
``` bash
//go.mod加入依赖
require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/zhuxiujia/GoMybatis v6.5.9+incompatible
)
```



实际使用mapper 定义xml内容，建议以*Mapper.xml文件存于项目目录中,在编辑xml时就可享受GoLand等IDE渲染和智能提示。生产环境可以使用statikFS把xml文件打包进程序里</br>
* main.go
``` xml
var xmlBytes = []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <select id="SelectAll">
        select * from biz_activity where delete_flag=1 order by create_time desc
    </select>
</mapper>
`)
```
``` go
import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //选择需要的数据库驱动导入
	"github.com/zhuxiujia/GoMybatis"
)
type ExampleActivityMapperImpl struct {
     SelectAll  func() ([]Activity, error)
}

func main() {
    var engine = GoMybatis.GoMybatisEngine{}.New()
	//Mysql链接格式 用户名:密码@(数据库链接地址:端口)/数据库名称,如root:123456@(***.com:3306)/test
	err := engine.Open("mysql", "*?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
	   panic(err)
	}
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	
	//加载xml实现逻辑到ExampleActivityMapperImpl
	engine.WriteMapperPtr(&exampleActivityMapperImpl, xmlBytes)

	//使用mapper
	result, err := exampleActivityMapperImpl.SelectAll(&result)
        if err != nil {
	   panic(err)
	}
	fmt.Println(result)
}
```
## 功能：模板标签CRUD 简化（依赖resultMap，同时带有 乐观锁.逻辑删除）
``` xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable 逻辑删除字段-->
    <!--logic_deleted 逻辑删除已删除字段-->
    <!--logic_undelete 逻辑删除 未删除字段-->
    <!--version_enable 乐观锁版本字段,支持int,int8,int16,int32,int64-->
    <resultMap id="BaseResultMap" tables="biz_activity">
        <id column="id" langType="string"/>
        <result column="name" langType="string"/>
        <result column="pc_link" langType="string"/>
        <result column="h5_link" langType="string"/>
        <result column="remark" langType="string"/>
        <result column="sort" langType="int"/>
        <result column="status" langType="status"/>
        <result column="version" langType="int"
                version_enable="true"/>
        <result column="create_time" langType="time.Time"/>
        <result column="delete_flag" langType="int"
                logic_enable="true"
                logic_undelete="1"
                logic_deleted="0"/>
    </resultMap>

    <!--模板标签: columns wheres sets 支持逗号,分隔表达式，*?* 为判空表达式-->

    <!--插入模板:默认id="insertTemplate,test="field != null",where自动设置逻辑删除字段,支持批量插入" -->
    <insertTemplate/>
    <!--查询模板:默认id="selectTemplate,where自动设置逻辑删除字段-->
    <selectTemplate wheres="name?name = #{name}"/>
    <!--更新模板:默认id="updateTemplate,set自动设置乐观锁版本号-->
    <updateTemplate sets="name?name = #{name},remark?remark=#{remark}" wheres="id?id = #{id}"/>
    <!--删除模板:默认id="deleteTemplate,where自动设置逻辑删除字段-->
    <deleteTemplate wheres="name?name = #{name}"/>
</mapper>    
```
xml对应以下定义的Mapper结构体方法,然后将生成对应的SQL语句
旧版本使用mapperParams的tag，单词太长容易拼错，新版改为args
```go
type Activity struct {
	Id         string    `json:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	PcLink     string    `json:"pcLink"`
	H5Link     string    `json:"h5Link"`
	Remark     string    `json:"remark"`
	Version    int       `json:"version"`
	CreateTime time.Time `json:"createTime"`
	DeleteFlag int       `json:"deleteFlag"`
}
type ExampleActivityMapper struct {
    //调用即可生成sql(带有逻辑删除)  select * from biz_activity where delete_flag = 1 and name = #{name}
	SelectTemplate      func(name string) ([]Activity, error) `args:"name"`
	InsertTemplate      func(arg Activity) (int64, error)
	InsertTemplateBatch func(args []Activity) (int64, error) `args:"args"`
    //生成sql(带有乐观锁.逻辑删除)  update biz_activity set name = #{name},remark=#{remark},version=#{version+1} where delete_flag = 1 and id = #{id} and version = #{version}
	UpdateTemplate      func(arg Activity) (int64, error)    `args:"name"`
	DeleteTemplate      func(name string) (int64, error)     `args:"name"`
}
```

## 功能：动态数据源
``` go
        //添加第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
    var engine = GoMybatis.GoMybatisEngine{}.New()
	engine.Open("mysql", MysqlUri)//添加第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
	var router = GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
		//根据包名路由指向数据源
		if strings.Contains(mapperName, "example.") {
			var url = MysqlUri//第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
			fmt.Println(url)
			return &url
		}
		return nil
	})
	engine.SetDataSourceRouter(&router)
```
## 功能：自定义日志输出
``` go
	engine.SetLogEnable(true)
	engine.SetLog(&GoMybatis.LogStandard{
		PrintlnFunc: func(messages []byte) {
		  //do someting save messages
		},
	})
```
## 功能：异步日志接口（可自定义日志输出）
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/log_system.png)

 ## 功能：嵌套事务传播处理器
 <table>
 <thead>
 <tr><th>事务类型</th>
 <th>说明</th>
 </tr>
 </thead>
 <tbody><tr><td>PROPAGATION_REQUIRED</td><td>表示如果当前事务存在，则支持当前事务。否则，会启动一个新的事务。默认事务类型。</td></tr>
 <tr><td>PROPAGATION_SUPPORTS</td><td>表示如果当前事务存在，则支持当前事务，如果当前没有事务，就以非事务方式执行。</td></tr>
 <tr><td>PROPAGATION_MANDATORY</td><td>表示如果当前事务存在，则支持当前事务，如果当前没有事务，则返回事务嵌套错误。</td></tr>
 <tr><td>PROPAGATION_REQUIRES_NEW</td><td>表示新建一个全新Session开启一个全新事务，如果当前存在事务，则把当前事务挂起。</td></tr>
 <tr><td>PROPAGATION_NOT_SUPPORTED</td><td>表示以非事务方式执行操作，如果当前存在事务，则新建一个Session以非事务方式执行操作，把当前事务挂起。</td></tr>
 <tr><td>PROPAGATION_NEVER</td><td>表示以非事务方式执行操作，如果当前存在事务，则返回事务嵌套错误。</td></tr>
 <tr><td>PROPAGATION_NESTED</td><td>表示如果当前事务存在，则在嵌套事务内执行，如嵌套事务回滚，则只会在嵌套事务内回滚，不会影响当前事务。如果当前没有事务，则进行与PROPAGATION_REQUIRED类似的操作。</td></tr>
 <tr><td>PROPAGATION_NOT_REQUIRED</td><td>表示如果当前没有事务，就新建一个事务,否则返回错误。</td></tr></tbody>
 </table>
 
 ``` go
 //嵌套事务的服务
type TestService struct {
	exampleActivityMapper *ExampleActivityMapper //服务包含一个mapper操作数据库，类似java spring mvc
	UpdateName   func(id string, name string) error   `tx:"" rollback:"error"`
	UpdateRemark func(id string, remark string) error `tx:"" rollback:"error"`
}
func main()  {
	var testService TestService
	testService = TestService{
		exampleActivityMapper: &exampleActivityMapper,
		UpdateRemark: func(id string, remark string) error {
			testService.exampleActivityMapper.SelectByIds([]string{id})
			panic(errors.New("业务异常")) // panic 触发事务回滚策略
			return nil                   // rollback:"error"指定了返回error类型 且不为nil 就会触发事务回滚策略
		},
		UpdateName: func(id string, name string) error {
			testService.exampleActivityMapper.SelectByIds([]string{id})
			return nil
		},
	}
	GoMybatis.AopProxyService(&testService, &engine)//必须使用AOP代理服务的func
	testService.UpdateRemark("1","remark")
}
```
 
 
 
 
 
  ## 功能：XML/Mapper生成器- 根据struct结构体生成*mapper.xml
``` go
  //step1 定义你的数据库模型,必须包含 json注解（默认为数据库字段）, gm:""注解指定 值是否为 id,version乐观锁,logic逻辑软删除
  type UserAddress struct {
	Id            string `json:"id" gm:"id"`
	UserId        string `json:"user_id"`
	RealName      string `json:"real_name"`
	Phone         string `json:"phone"`
	AddressDetail string `json:"address_detail"`

	Version    int       `json:"version" gm:"version"`
	CreateTime time.Time `json:"create_time"`
	DeleteFlag int       `json:"delete_flag" gm:"logic"`
}
```
* 第二步，在你项目main 目录下建立一个 XmlCreateTool.go 内容如下
``` go
func main() {
	var bean = UserAddress{} //此处只是举例，应该替换为你自己的数据库模型
	GoMybatis.OutPutXml(reflect.TypeOf(bean).Name()+"Mapper.xml", GoMybatis.CreateXml("biz_"+GoMybatis.StructToSnakeString(bean), bean))
}
```
* 第三步，执行命令，在当前目录下得到 UserAddressMapper.xml文件
``` go
go run XmlCreateTool.go
```
* 以下是自动生成的xml文件内容
``` xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable 逻辑删除字段-->
    <!--logic_deleted 逻辑删除已删除字段-->
    <!--logic_undelete 逻辑删除 未删除字段-->
    <!--version_enable 乐观锁版本字段,支持int,int8,int16,int32,int64-->
    <resultMap id="BaseResultMap" tables="biz_user_address">
    <id column="id" property="id"/>
	<result column="id" property="id" langType="string"   />
	<result column="user_id" property="user_id" langType="string"   />
	<result column="real_name" property="real_name" langType="string"   />
	<result column="phone" property="phone" langType="string"   />
	<result column="address_detail" property="address_detail" langType="string"   />
	<result column="version" property="version" langType="int" version_enable="true"  />
	<result column="create_time" property="create_time" langType="Time"   />
	<result column="delete_flag" property="delete_flag" langType="int"  logic_enable="true" logic_undelete="1" logic_deleted="0" />
    </resultMap>
</mapper>
```
 
 
## 为了让您快速学习此框架，建议查看在实际项目中的应用 [点击-GoMybatis商城](https://github.com/zhuxiujia/GoMybatisMall) 
 
 
## 建议使用的框架（已应用在生产环境）配合GoMybatis
#### [easy_mvc](https://github.com/zhuxiujia/easy_mvc) 
* 整体基于反射tag，所有配置（包括http方法，路径，参数，swagger文档参数）都集中于你定义的函数之后
* 轻量 完全兼容标准库的http，意味着和标准库一般稳定，可以混合搭配使用，扩展性极高
* 拦截器 支持（例如非常方便的检查用户登录，提取用户登录数据，支持JWT token，Oath2Token更加方便的接入）
* 过滤器 支持
* 全局错误处理器链 支持
* 使用tag 定义 http请求参数，包含 *int,*string,*float 同时支持标准库的 writer http.ResponseWriter, request *http.Request
* Json参数支持（app端上传时需要Header，Content-Type设置为application/json）
* 支持参数默认值 只需在tag中 定义，例如 func(phone string, pwd string, age *int) interface{} arg:"phone,pwd,age:1"  其中 arg没有传参则默认为1
* 指针参数可为空（nil）非指针参数 如果没有值框架会拦截
* root path支持，类似spring controller定义一个基础的path加控制器具体方法的http path
* 支持swagger ui 动态文档，免生成任何中间go文件 基于Tag和反射实现的swagger动态文档
#### [easy_rpc](https://github.com/zhuxiujia/easyrpc)  （RPC框架，和GoMybatis配合更容易）
* 基于标准库rpc库修改而来,稳定,高性能,扩展性好
* 标准库默认使用func (* Type)Method(arg,*result) error 的模式,EasyRpc 则把方法移动到结构体里（方便动态代理和Aop以及各种扩展和定制）
* easyrpc同时支持 无参数，无返回值，或只有参数，只有返回值
* 支持注册defer函数  easyrpc.RegisterDefer(v,deferFunc) ，防止服务因为不可预知 painc 问题导致程序退出。defer函数可处理问题然后把错误发送还给客户端
#### [easyrpc_discovery](https://github.com/zhuxiujia/easyrpc_discovery)  服务发现
* 自带负载均衡算法 随机 加权轮询 源地址哈希法
* 基于easyrpc,类似标准库的api，定义服务没有标准库的要求那么严格（可选不传参数，或者只有一个参数，只有一个返回值） https://github.com/zhuxiujia/easyrpc
* 基于easyrpc，负载均衡算法，失败重试，支持动态代理，支持GoMybatis事务，AOP代理，事务嵌套，tag定义事务
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/easy_consul.png)




## 为了让您快速学习此框架，建议查看在实际项目中的应用 [点击-GoMybatis商城](https://github.com/zhuxiujia/GoMybatisMall) 



## 新特性规划
* 内置方法explain目录下所有sql，获取每条sql是否走索引.避免上线后因为sql扫全表导致卡顿和死锁悲剧
* 新项目 https://github.com/rbatis/rbatis 提供高性能，无GC,无并发安全问题，内存安全的rust语言 orm框架

## 请及时关注版本，尽可能使用最新版本(稳定，修复bug) 
* 不管是商业用途还是个人使用GoMybatis项目，必须在Issues里留言您的项目名称+联系方式 ！

## 联系方式：微信号 zxj347284221

## 欢迎右上角star 或捐赠赞助~
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/wx_account.jpg)

