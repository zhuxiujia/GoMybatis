
# id package,推荐使用雪花算法,以获取分布式的id并且是稠密索引id

* snowflake clone from github   https://github.com/bwmarrin/snowflake



### how to use it

```go
  n, _ := NewNode(0)
	id := n.Generate()
	fmt.Println(id)
```
