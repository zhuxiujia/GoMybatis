# GoMybatis
GoMybatis 是根据java版 Mybatis3 的实现,基于Gorm和govaluate表达式及反射实现。
GoMybatis 内部在初始化时反射分析mapper xml生成golang的func代码，默认支持绝大部分的Java版的mybatis标签和规范,
支持标签
`<select>,<update>,<insert>,<delete>,<trim>,<if>,<set>,<foreach>`

