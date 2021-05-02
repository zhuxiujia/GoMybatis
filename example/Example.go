package example

import "time"

//mysql链接格式为         用户名:密码@(数据库链接地址:端口)/数据库名称   例如root:root@(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local
//此处请按格式填写你的mysql链接，这里用*号代替
const MysqlUri = "*"

//定义数据库模型
//例子：Activity 活动数据,注意GoMybatis json tag 等价于数据库字段
type Activity struct {
	Id         string    `json:"id,omitempty"`
	Uuid       string    `json:"uuid"`
	Name       string    `json:"name"`
	PcLink     string    `json:"pc_link"`
	H5Link     string    `json:"h5_link"`
	Remark     string    `json:"remark"`
	Sort       int       `json:"sort"`
	Status     int       `json:"status"`
	Version    int       `json:"version"`
	CreateTime time.Time `json:"create_time"`
	DeleteFlag int       `json:"delete_flag"`
}
