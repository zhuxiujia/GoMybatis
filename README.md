# GoMybatis
![Image text](https://github.com/zhuxiujia/gomybatis.io/assets/vuetify.png)
# 网站/文档地址
https://zhuxiujia.github.io/gomybatis.io/
# 优势
GoMybatis 是根据java版 Mybatis3 的实现,基于Go标准库和github.com/Knetic/govaluate表达式及github.com/beevik/etree读取Xml解析,github.com/satori/go.uuid生成库 实现。
GoMybatis 内部在初始化时反射分析mapper xml生成golang的func代码，默认支持绝大部分的Java版的mybatis标签和规范,
### 已支持绝大部分标签
`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<foreach><resultMap>`
### 已支持本地和远程事务,方便处于 单数据库(Mysql,postgresql)-分布式数据库（TiDB,cockroachdb...）过渡期间的微服务
![Image text](https://github.com/zhuxiujia/gomybatis.io/raw/master/docs/tx.png)

## 使用教程,
##### 代码文件请移步Github文件夹https://github.com/zhuxiujia/GoMybatis/tree/master/example
各种数据库驱动支持
```
 Mysql: github.com/go-sql-driver/mysql
 MyMysql: github.com/ziutek/mymysql/godrv
 Postgres: github.com/lib/pq
 Tidb: github.com/pingcap/tidb
 SQLite: github.com/mattn/go-sqlite3
 MsSql: github.com/denisenkom/go-mssqldb
 MsSql: github.com/lunny/godbc
 Oracle: github.com/mattn/go-oci8
 ```
设置好GoPath,用go get 命令下载GoMybatis和对应的数据库驱动
```
go get github.com/zhuxiujia/GoMybatis
go get github.com/go-sql-driver/mysql
```
mapper.go 文件案例
```
//定义mapper文件的接口和结构体，也可以只定义结构体就行
//mapper.go文件 函数参数（自定义结构体参数（属性必须大写），为指针类型的返回数据,*GoMybatis.Session作为该sql执行的session） error 为返回错误
type ExampleActivityMapperImpl struct {
	SelectAll         func(result *[]Activity) error
	SelectByCondition func(name string, startTime time.Time, endTime time.Time, page int, size int, result *[]Activity) error `mapperParams:"name,startTime,endTime,page,size"`
	UpdateById        func(session *GoMybatis.Session, arg Activity, result *int64) error //只要参数中包含有*GoMybatis.Session的类型，框架默认使用传入的session对象，用于自定义事务
	Insert            func(arg Activity, result *int64) error
	CountByCondition  func(name string, startTime time.Time, endTime time.Time, result *int) error                            `mapperParams:"name,startTime,endTime"`
}
```

xml文件案例:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "https://github.com/zhuxiujia/GoMybatis/blob/master/mybatis-3-mapper.dtd">
<mapper namespace="ActivityMapperImpl">
    <resultMap id="BaseResultMap" type="example.Activity">
        <id column="id" property="id" goType="string"/>
        <result column="name" property="name" goType="string"/>
        <result column="pc_link" property="pcLink" goType="string"/>
        <result column="h5_link" property="h5Link" goType="string"/>
        <result column="remark" property="remark" goType="string"/>
        <result column="create_time" property="createTime" goType="time.Time"/>
        <result column="delete_flag" property="deleteFlag" goType="int"/>
    </resultMap>
    <!--SelectAll(result *[]Activity)error-->
    <select id="selectAll" resultMap="BaseResultMap">
        select * from biz_activity where delete_flag=1 order by create_time desc
    </select>
</mapper>
```
实际使用mapper
```
import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/zhuxiujia/GoMybatis"
)
func main() {
  var err error
  	//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:123456@(***.mysql.rds.aliyuncs.com:3306)/test
  	engine, err := GoMybatis.Open("mysql", "*?charset=utf8&parseTime=True&loc=Local") //此处请按格式填写你的mysql链接，这里用*号代替
  	if err != nil {
  		panic(err.Error())
  	}
  
  	file, err := os.Open("Example_ActivityMapper.xml")
  	if err != nil {
  		panic(err)
  	}
  	defer file.Close()
  
  	bytes, _ := ioutil.ReadAll(file)
  	var exampleActivityMapperImpl ExampleActivityMapperImpl
  	//设置对应的mapper xml文件
  	GoMybatis.UseProxyMapperByEngine(&exampleActivityMapperImpl, bytes, engine)
  
  	//使用mapper
  	var result []Activity
  	exampleActivityMapperImpl.SelectAll(&result)
  
  	fmt.Println(result)
}
```

## TODO 期待功能
-`<sql><include>` 标签支持（进行中）</br>
-针对于 GoLand 的xml生成插件,可以使用鼠标右键点击一键生成CRUD基础XML(待支持..)</br>

