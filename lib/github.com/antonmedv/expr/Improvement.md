
## 此次列出改进项

## 1 改进支持表达式取指针属性值会直接取到实际值
 设foo.Bar = *int = 1</br>
 原版 表达式 ' foo.Bar ' 执行eval后 结果为  0x123ffcc</br>
 改进 表达式 ' foo.Bar ' 执行eval后 结果为  1</br>
eval.go 中 加入取实际值代码</br>
```
//指针转换真实值 在 func (n binaryNode) Eval(env interface{}) (interface{}, error)
	var reflectLeft = reflect.ValueOf(left)
	if reflectLeft.IsValid() && (reflectLeft.Type().Kind() == reflect.Interface || reflectLeft.Type().Kind() == reflect.Ptr) {
		reflectLeft = GetDeepPtr(reflectLeft)
		if reflectLeft.IsValid() && reflectLeft.CanInterface(){
			left = reflectLeft.Interface()
		}
	}
```
```
//在func Run(node Node, env interface{}) (out interface{}, err error) 
var resultV = reflect.ValueOf(result)
	if resultV.IsValid() && (resultV.Type().Kind() == reflect.Interface || resultV.Type().Kind() == reflect.Ptr) {
		resultV = GetDeepPtr(resultV)
		if resultV.IsValid() && resultV.CanInterface(){
			result = resultV.Interface()
			return result, nil
		}
	}
```
### 2 改进同时支持 表达式  ' foo.Bar == nil '  ' foo.Bar == null '
```
//在 parser.go 中
case "null":
			return nilNode{}, nil
```
### 3 改进支持字符串相加操作符为+ 号
 原版表达式 "'foo.Bar'~ 'foo.Bar'"" 执行eval后 结果为  'foo.Barfoo.Bar'</br>
 改进表达式 "'foo.Bar'+ 'foo.Bar'"" 执行eval后 结果为  'foo.Barfoo.Bar'</br>