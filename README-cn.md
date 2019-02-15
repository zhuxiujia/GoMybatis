# Go语言sql数据库orm框架
* 中文
* [English](README.md) 
[![Go Report Card](https://goreportcard.com/badge/github.com/zhuxiujia/GoMybatis)](https://goreportcard.com/report/github.com/zhuxiujia/GoMybatis)
[![Build Status](https://travis-ci.com/zhuxiujia/GoMybatis.svg?branch=master)](https://travis-ci.com/zhuxiujia/GoMybatis)
[![GoDoc](https://godoc.org/github.com/zhuxiujia/GoMybatis?status.svg)](https://godoc.org/github.com/zhuxiujia/GoMybatis)
[![Coverage Status](https://coveralls.io/repos/github/zhuxiujia/GoMybatis/badge.svg?branch=master)](https://coveralls.io/github/zhuxiujia/GoMybatis?branch=master)
[![codecov](https://codecov.io/gh/zhuxiujia/GoMybatis/branch/master/graph/badge.svg)](https://codecov.io/gh/zhuxiujia/GoMybatis)


![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/vuetify.png)
### 网站 https://zhuxiujia.github.io/gomybatis.io/info.html
# 优势
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-协程高并发</a>，假设数据库响应时间为0，在6核16GB PC上框架每秒事务数可达246982Tps/s,耗时仅仅0.4s<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-事务支持</a>，session灵活插拔，兼容过渡期微服务<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-动态SQL</a>，在xml中可灵活运用`<if>`判断，`<foreach>`数组/map，`<resultMap>,<bind>`等等java框架Mybatis包含的15种实用功能<br>
`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<where>,<foreach>,<resultMap>,<bind>,<choose><when><otherwise>,<sql><include>`<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-多数据库支持Mysql,Postgres,Tidb,SQLite,Oracle....等等更多</a><br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-快速上手</a>基于反射动态代理,无需go generate生成*.go等中间代码，xml读取后可直接调用函数<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-表达式和xml面向对象设计</a>假如foo.Bar 这个属性是指针,那么在xml中调用 foo.Bar 则会取实际值,完全避免使用&和*符号指针操作<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-支持动态多数据源</a>可以使用路由engine.SetDataSourceRouter自定义多数据源规则<br>
<a href="https://zhuxiujia.github.io/gomybatis.io/info.html">-模板标签</a>高开发效率的模板，一行代码实现增删改查，逻辑删除，乐观锁版本号<br>
#### 异步消息队列日志系统
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/log_system.png)
#### 通过远程代替微服务成员 处理事务支持  处于 单数据库(Mysql,postgresql)-分布式数据库（TiDB,cockroachdb...）过渡期间的微服务
![Image text](https://zhuxiujia.github.io/gomybatis.io/assets/tx.png)

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
	"github.com/zhuxiujia/GoMybatis"
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


#### v2019.1.19 新增了gojee引擎（改进了原作者源码取消了开头的"."符号,例如.a.b变成 a.b）和expr表达式引擎（改进了原作者源码指针的bug,加入字符串相加操作例如'a'+'b'）
https://github.com/zhuxiujia/GoMybatis/tree/master/lib/github.com/antonmedv/expr
https://github.com/zhuxiujia/GoMybatis/tree/master/lib/github.com/nytlabs/gojee
https://github.com/zhuxiujia/GoMybatis/tree/master/lib/github.com/Knetic/govaluate
<table border="1">
     <tr>
        <td>表达式引擎</td>
        <td>是否支持指针参数/null/nil</td>
        <td>执行效率</td>
	<td>and or 命令性能损耗</td>
        <td>表达式功能</td>
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
为了执行效率 框架默认使用 github.com/antonmedv/expr作为默认选项，你也可以自定义调用GoMybatis.WriteMapper()参数中SqlBuilder的参数自行选择加入ExpressionEngine

## 请及时关注版本，及时升级版本(新的功能，bug修复)
## TODO 期待功能路线（预览特性,有可能会更改）
-专属简化模板标签,媲美语言层orm框架的开发效率，和java端的mybatis啰嗦的语法告别（进行中）
-逻辑删除
-乐观锁
## 喜欢的老铁欢迎在右上角点下 star 关注和支持我们哈
