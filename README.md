# SQL mapper ORM framework for Golang
* English
* [中文](README-ch.md)   


[![Go Report Card](https://goreportcard.com/badge/github.com/zhuxiujia/GoMybatis)](https://goreportcard.com/report/github.com/zhuxiujia/GoMybatis)
[![Build Status](https://travis-ci.com/zhuxiujia/GoMybatis.svg?branch=master)](https://travis-ci.com/zhuxiujia/GoMybatis)
[![GoDoc](https://godoc.org/github.com/zhuxiujia/GoMybatis?status.svg)](https://godoc.org/github.com/zhuxiujia/GoMybatis)
[![Coverage Status](https://coveralls.io/repos/github/zhuxiujia/GoMybatis/badge.svg?branch=master)](https://coveralls.io/github/zhuxiujia/GoMybatis?branch=master)
[![codecov](https://codecov.io/gh/zhuxiujia/GoMybatis/branch/master/graph/badge.svg)](https://codecov.io/gh/zhuxiujia/GoMybatis)


![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/vuetify.png)
### Please read the documentation website carefully when using the tutorial. [DOC](https://zhuxiujia.github.io/gomybatis.io/#/getting-started)
# Powerful Features
* <a href="https://zhuxiujia.github.io/gomybatis.io/">High Performance</a>，The maximum number of transactions per second of a single computer can reach 751020 Tps/s, and the total time consumed is 0.14s (test environment returns simulated SQL data, concurrently 1000, total 100000, 6-core 16GB win10)<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Painless migration from Java to go</a>，Compatible with most Java(Mybatis3,Mybatis Plus) ，Painless migration of XML SQL files from Java Spring Mybatis to Go language（Modify only the javaType of resultMap to specify go language type for langType）<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Declarative transaction/AOP transaction/transaction Behavior</a>Only one line Tag is needed to define AOP transactions and transaction propagation behavior<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Extensible Log Interface</a>Asynchronous message queue day, SQL log in framework uses cached channel to realize asynchronous message queue logging<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">dynamic sql</a>，contains 15 utilities Features`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<where>,<foreach>,<resultMap>,<bind>,<choose><when><otherwise>,<sql><include>`<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Intelligent expression</a>Processing dynamic judgment and computation tasks（such as：`#{foo.Bar}#{arg+1}#{arg*1}#{arg/1}#{arg-1}`）,For example, write fuzzy queries `select * from table where phone like #{phone+'%'}`(Note the post-percentile query run in index)<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Dynamic Data Source</a>Multiple data sources can be customized to dynamically switch multiple database instances<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Template label（new）</a>One line of code to achieve add, delete, modify, delete logic, optimistic lock, but also retain perfect scalability (tag body can continue to expand SQL logic)<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Optimistic Lock（new）</a>`<updateTemplete>`Optimistic locks to prevent concurrent competition to modify records as much as possible<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">Logical deletion（new）</a>`<insertTemplete><updateTemplete><deleteTemplete><selectTemplete>`Logical deletion, prevent accidental deletion of data, data recovery is simple<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/">RPC/MVC Component Support（new）</a>To make the service perfect for RPC (reducing parameter restrictions), dynamic proxy, transaction subscription, easy integration and extension of micro services, click on the link https://github.com/zhuxiujia/easyrpc<br>


## Database Driver support table
``` bash
 //Traditional database
 Mysql:                             github.com/go-sql-driver/mysql
 MyMysql:                           github.com/ziutek/mymysql/godrv
 Postgres:                          github.com/lib/pq
 SQLite:                            github.com/mattn/go-sqlite3
 MsSql:                             github.com/denisenkom/go-mssqldb
 MsSql:                             github.com/lunny/godbc
 Oracle:                            github.com/mattn/go-oci8
 //Distributed NewSql database
 Tidb:                              github.com/pingcap/tidb
 CockroachDB:                       github.com/lib/pq
 ```
 
## Use tutorials
> Tutorial source code  https://github.com/zhuxiujia/GoMybatis/tree/master/example

Set up GoPath and download GoMybatis and the corresponding database driver with the go get command
``` bash
go get github.com/zhuxiujia/GoMybatis
go get github.com/go-sql-driver/mysql
```
In practice, we use mapper to define the content of xml. It is suggested that the * Mapper. XML file be stored in the project directory. When editing xml, we can enjoy IDE rendering and intelligent prompts such as GoLand. Production environments can use statikFS to package XML files in the process</br>
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
	_ "github.com/go-sql-driver/mysql" //Select the required database-driven imports
	"github.com/zhuxiujia/GoMybatis"
)
type ExampleActivityMapperImpl struct {
     SelectAll  func() ([]Activity, error)
}

func main() {
    var engine = GoMybatis.GoMybatisEngine{}.New()
	//Mysql link format user name: password @ (database link address: port)/database name, such as root: 123456 @(***.com: 3306)/test
	err := engine.Open("mysql", "*?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
	   panic(err)
	}
	var exampleActivityMapperImpl ExampleActivityMapperImpl
	
	//Loading XML implementation logic to ExampleActivity Mapper Impl
	engine.WriteMapperPtr(&exampleActivityMapperImpl, xmlBytes)

	//use mapper
	result, err := exampleActivityMapperImpl.SelectAll(&result)
        if err != nil {
	   panic(err)
	}
	fmt.Println(result)
}
```
## Features: Template tag CRUD simplification (must rely on a resultMap tag)
``` xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable Logical Delete Fields-->
    <!--logic_deleted Logically delete deleted fields-->
    <!--logic_undelete Logically Delete Undeleted Fields-->
    <!--version_enable Optimistic lock version field, support int, int8, int16, int32, Int64-->
    <resultMap id="BaseResultMap" tables="biz_activity">
        <id column="id" property="id"/>
        <result column="name" property="name" langType="string"/>
        <result column="pc_link" property="pcLink" langType="string"/>
        <result column="h5_link" property="h5Link" langType="string"/>
        <result column="remark" property="remark" langType="string"/>
        <result column="version" property="version" langType="int"
                version_enable="true"/>
        <result column="create_time" property="createTime" langType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" langType="int"
                logic_enable="true"
                logic_undelete="1"
                logic_deleted="0"/>
    </resultMap>
    <!--Template tags: columns wheres sets support commas, separating expressions, *?* as null expressions-->
    <!--Insert Template: Default id="insertTemplete,test="field != null",where Automatically set logical deletion fields to support batch insertion" -->
    <insertTemplete/>
    <!--Query template: default id="selectTemplete,where Automatically Set Logical Delete Fields-->
    <selectTemplete wheres="name?name = #{name}"/>
    <!-- Update template: default id="updateTemplete,set Automatically Setting Optimistic Lock Version Number-->
    <updateTemplete sets="name?name = #{name},remark?remark=#{remark}" wheres="id?id = #{id}"/>
    <!--Delete template: default id="deleteTemplete,where Automatically Set Logical Delete Fields-->
    <deleteTemplete wheres="name?name = #{name}"/>
</mapper>    
```
XML corresponds to the Mapper structure method defined below
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
	SelectTemplete      func(name string) ([]Activity, error) `mapperParams:"name"`
	InsertTemplete      func(arg Activity) (int64, error)
	InsertTempleteBatch func(args []Activity) (int64, error) `mapperParams:"args"`
	UpdateTemplete      func(arg Activity) (int64, error)    `mapperParams:"name"`
	DeleteTemplete      func(name string) (int64, error)     `mapperParams:"name"`
}
```

## Features：Dynamic Data Source
``` go
        //To add a second MySQL database, change Mysql Uri to your second data source link
	GoMybatis.Open("mysql", MysqlUri)
	//Dynamic Data Source Routing
	var router = GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
		//Point to the data source according to the packet name routing
		if strings.Contains(mapperName, "example.") {
			var url = MysqlUri//The second MySQL database, please change Mysql Uri to your second data source link
			fmt.Println(url)
			return &url
		}
		return nil
	})
```
## Features：Custom log output
``` go
	engine.SetLogEnable(true)
	engine.SetLog(&GoMybatis.LogStandard{
		PrintlnFunc: func(messages []byte) {
		  //do someting save messages
		},
	})
```
## Features：Asynchronous log interface (customizable log output)
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/log_system.png)

 ## Features：Transaction Propagation Processor (Nested Transactions)
 <table>
 <thead>
 <tr><th>Transaction type</th>
 <th>Explain</th>
 </tr>
 </thead>
 <tbody><tr><td>PROPAGATION_REQUIRED</td><td>Represents that if the current transaction exists, the current transaction is supported. Otherwise, a new transaction will be started. Default transaction type.</td></tr>
 <tr><td>PROPAGATION_SUPPORTS</td><td>Represents that if the current transaction exists, the current transaction is supported, and if there is no transaction at present, it is executed in a non-transactional manner.</td></tr>
 <tr><td>PROPAGATION_MANDATORY</td><td>Represents that if the current transaction exists, the current transaction is supported, and if no transaction exists, the transaction nesting error is returned.</td></tr>
 <tr><td>PROPAGATION_REQUIRES_NEW</td><td>Represents that a new Session opens a new transaction and suspends the current transaction if it currently exists.</td></tr>
 <tr><td>PROPAGATION_NOT_SUPPORTED</td><td>Represents that an operation is performed in a non-transactional manner. If a transaction exists, a new Session is created to perform the operation in a non-transactional manner, suspending the current transaction.</td></tr>
 <tr><td>PROPAGATION_NEVER</td><td>Represents that an operation is executed in a non-transactional manner and returns a transaction nesting error if a transaction currently exists.</td></tr>
 <tr><td>PROPAGATION_NESTED</td><td>Represents that if the current transaction exists, it will be executed within the nested transaction. If the nested transaction rolls back, it will only roll back within the nested transaction and will not affect the current transaction. If there is no transaction at the moment, do something similar to PROPAGATION_REQUIRED.</td></tr>
 <tr><td>PROPAGATION_NOT_REQUIRED</td><td>Represents that if there is currently no transaction, a new transaction will be created, otherwise an error will be returned.</td></tr></tbody>
 </table>
 
 ``` go
 //Nested transaction services
type TestService struct {
	exampleActivityMapper *ExampleActivityMapper //The service contains a mapper operation database similar to Java spring MVC
	UpdateName   func(id string, name string) error   `tx:"" rollback:"error"`
	UpdateRemark func(id string, remark string) error `tx:"" rollback:"error"`
}
func main()  {
	var testService TestService
	testService = TestService{
		exampleActivityMapper: &exampleActivityMapper,
		UpdateRemark: func(id string, remark string) error {
			testService.exampleActivityMapper.SelectByIds([]string{id})
			panic(errors.New("Business exceptions")) // panic Triggered transaction rollback strategy
			return nil                   // rollback:"error" A transaction rollback policy is triggered if the error type is returned and is not nil
		},
		UpdateName: func(id string, name string) error {
			testService.exampleActivityMapper.SelectByIds([]string{id})
			return nil
		},
	}
	GoMybatis.AopProxyService(&testService, &engine)//Func must use AOP proxy service
	testService.UpdateRemark("1","remark")
}
```
 
 
 
 
 
  ## Features：XML/Mapper Generator - Generate * mapper. XML from struct structure
``` go
  //step1 To define your database model, you must include JSON annotations (default database fields), gm:"" annotations specifying whether the value is id, version optimistic locks, and logic logic soft deletion.
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
* Step 2: Create an Xml CreateTool. go in the main directory of your project as follows
```
func main() {
	var bean = UserAddress{} //Here's just an example, which should be replaced by your own database model
	GoMybatis.OutPutXml(reflect.TypeOf(bean).Name()+"Mapper.xml", GoMybatis.CreateXml("biz_"+GoMybatis.StructToSnakeString(bean), bean))
}
```
* Third, execute the command to get the UserAddressMapper. XML file in the current directory
``` go
go run XmlCreateTool.go
```
* The following is the content of the automatically generated XML file
``` xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper>
    <!--logic_enable Logical Delete Fields-->
    <!--logic_deleted Logically delete deleted fields-->
    <!--logic_undelete Logically Delete Undeleted Fields-->
    <!--version_enable Optimistic lock version field, support int, int8, int16, int32, Int64-->
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
 
 
 
 
 
 


## Components (RPC, JSONRPC, Consul) - With GoMybatis
* https://github.com/zhuxiujia/easy_mvc //mvc
* https://github.com/zhuxiujia/easyrpc  //easyrpc
* https://github.com/zhuxiujia/easyrpc_discovery  //easyrpc discovery
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/easy_consul.png)







## Please pay attention to the version in time, upgrade the version in time (new features, bug fix). For projects using GoMybatis, please leave your project name + contact information in Issues.

## Welcome to Star or Wechat Payment Sponsorship at the top right corner~
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/wx_account.jpg)
