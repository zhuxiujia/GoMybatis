package example

import "time"

//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:root@(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local
//此处请按格式填写你的mysql链接，这里用*号代替
const MysqlUri = "*"

//定义数据库模型
//例子：Activity 活动数据
type Activity struct {
	Id         string    `json:"id"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	PcLink     string    `json:"pcLink"`
	H5Link     string    `json:"h5Link"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
	DeleteFlag int       `json:"deleteFlag"`
}
