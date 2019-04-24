# Go语言sql数据库orm框架
* 中文
* [English](README-en.md) 

[![Go Report Card](https://goreportcard.com/badge/github.com/zhuxiujia/GoMybatis)](https://goreportcard.com/report/github.com/zhuxiujia/GoMybatis)
[![Build Status](https://travis-ci.com/zhuxiujia/GoMybatis.svg?branch=master)](https://travis-ci.com/zhuxiujia/GoMybatis)
[![GoDoc](https://godoc.org/github.com/zhuxiujia/GoMybatis?status.svg)](https://godoc.org/github.com/zhuxiujia/GoMybatis)
[![Coverage Status](https://coveralls.io/repos/github/zhuxiujia/GoMybatis/badge.svg?branch=master)](https://coveralls.io/github/zhuxiujia/GoMybatis?branch=master)
[![codecov](https://codecov.io/gh/zhuxiujia/GoMybatis/branch/master/graph/badge.svg)](https://codecov.io/gh/zhuxiujia/GoMybatis)


![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/vuetify.png)
### 使用教程请仔细阅读文档网站 https://zhuxiujia.github.io/gomybatis.io/info.html
# 优势
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">高性能</a>，单机每秒事务数最高可达456621Tps/s,总耗时0.22s （测试环境 返回模拟的sql数据，并发1000，总数100000，6核16GB win10）<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">事务</a>，session灵活插拔，兼容过渡期微服务<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">异步日志</a>异步消息队列日,框架内sql日志使用带缓存的channel实现 消息队列异步记录日志<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">动态SQL</a>，在xml中`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<where>,<foreach>,<resultMap>,<bind>,<choose><when><otherwise>,<sql><include>`等等java框架Mybatis包含的15种实用功能<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">多数据库</a>Mysql,Postgres,Tidb,SQLite,Oracle....等等更多<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">无依赖</a>基于反射动态代理,无需go generate生成*.go等中间代码，xml读取后可直接调用函数<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">智能表达式</a>`#{foo.Bar}``#{arg+1}``#{arg*1}``#{arg/1}``#{arg-1}`不但可以处理简单判断和计算任务，同时在取值时 假如foo.Bar 这个属性是指针,那么调用 foo.Bar 则会取指针指向的实际值,完全避免解引用操作<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">动态数据源</a>可以使用路由engine.SetDataSourceRouter自定义多数据源规则<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">模板标签</a>一行代码实现增删改查，逻辑删除，乐观锁（基于版本号更新）极大减轻CRUD操作的心智负担<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">乐观锁</a>`<updateTemplete>`支持通过修改版本号实现的乐观锁<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">逻辑删除</a>`<insertTemplete>``<updateTemplete>``<deleteTemplete>``<selectTemplete>`均支持逻辑删除<br>
* <a href="https://zhuxiujia.github.io/gomybatis.io/info.html">8种事务传播行为</a>复刻Spring MVC的事务传播行为功能<br>


## 使用教程

> 示例源码https://github.com/zhuxiujia/GoMybatis/tree/master/example

设置好GoPath,用go get 命令下载GoMybatis和对应的数据库驱动
```
go get github.com/zhuxiujia/GoMybatis
go get github.com/go-sql-driver/mysql
```
实际使用mapper
```
import (
	_ "github.com/go-sql-driver/mysql" //导入mysql驱动
	"GoMybatis"
	"fmt"
	"time"
)

//定义xml内容，建议以*Mapper.xml文件存于项目目录中,在编辑xml时就可享受GoLand等IDE渲染和智能提示。生产环境可以使用statikFS把xml文件打包进程序里

var xmlBytes = []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
"https://raw.githubusercontent.com/zhuxiujia/GoMybatis/master/mybatis-3-mapper.dtd">
<mapper namespace="ActivityMapperImpl">
    <!--SelectAll(result *[]Activity)error-->
    <select id="selectAll">
        select * from biz_activity where delete_flag=1 order by create_time desc
    </select>
</mapper>
`)

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
## 动态数据源
```
        //添加第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
	GoMybatis.Open("mysql", MysqlUri)
	//动态数据源路由
	var router = GoMybatis.GoMybatisDataSourceRouter{}.New(func(mapperName string) *string {
		//根据包名路由指向数据源
		if strings.Contains(mapperName, "example.") {
			var url = MysqlUri//第二个mysql数据库,请把MysqlUri改成你的第二个数据源链接
			fmt.Println(url)
			return &url
		}
		return nil
	})
```
## 自定义日志输出
```
	engine.SetLogEnable(true)
	engine.SetLog(&GoMybatis.LogStandard{
		PrintlnFunc: func(messages []byte) {
		},
	})
```
## 异步日志-基于消息队列日志
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/log_system.png)
数据库驱动列表
```
 Mysql: github.com/go-sql-driver/mysql
 MyMysql: github.com/ziutek/mymysql/godrv
 Postgres: github.com/lib/pq
 Tidb: github.com/pingcap/tidb
 SQLite: github.com/mattn/go-sqlite3
 MsSql: github.com/denisenkom/go-mssqldb
 MsSql: github.com/lunny/godbc
 Oracle: github.com/mattn/go-oci8
 CockroachDB(Postgres): github.com/lib/pq
 ```
 
 ## 嵌套事务-事务传播行为
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


## 多种表达式引擎可选（表达式引擎接口ExpressionEngine.go 负责表达式("foo != nil"...)的判断和取值）
<table border="1">
     <tr>
        <td>表达式引擎</td>
        <td>是否支持指针参数/null/nil</td>
        <td>执行效率</td>
	<td>and or 命令性能损耗</td>
        <td>表达式功能</td>
    </tr>
     <tr>
             <td>GoFastExpress</td>
             <td>支持null和nil和指针</td>
             <td>快-实测比expr快</td>
    	     <td>还行</td>
             <td>少</td>
        </tr>
    <tr>
         <td>expr</td>
         <td>支持null和nil和指针</td>
         <td>快-实测比govaluate快</td>
	     <td>还行</td>
         <td>一般</td>
    </tr>
    <tr>
          <td>gojee</td>
          <td>支持null和指针</td>
          <td>慢-每次检查表达式都有json序列化和反序列化操作</td>
	  <td>大</td>
          <td>多</td>
    </tr>
    <tr>
           <td>govaluate</td>
           <td>不支持null和nil和指针</td>
           <td>中等</td>
	   <td>一般</td>
           <td>一般</td>
    </tr>
</table>
为了执行效率 框架默认使用 github.com/zhuxiujia/GoFastExpress 表达式引擎作为默认选项，你也可以自定义调用GoMybatis.WriteMapper()参数中SqlBuilder的参数自行选择加入ExpressionEngine

## 请及时关注版本，及时升级版本(新的功能，bug修复)
## TODO 未来新特性（可能会更改）
* 模板标签,一行代码crud（已完成）
* 逻辑删除(已完成)
* 乐观锁（已完成）
* 重构SqlBuilder，使用抽象语法树代替递归，获得更好的维护性可读性(已完成)
* <a href="https://github.com/zhuxiujia/RustMybatis">Rust语言版本的 RustMybatis</a>基于Rust语言和LLVM编译器`不输于C++性能``无运行时``无GC``无`,预期会在 -并发数,-内存消耗, 都将远超go语言版 (开发中) 敬请期待~
## 喜欢的老铁欢迎在右上角点下 star 关注和支持我们哈


