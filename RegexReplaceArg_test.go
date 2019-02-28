package GoMybatis

import (
	"fmt"
	"testing"
)

type TestBean struct {
	Name  string
	Child TestBeanChild
}
type TestBeanChild struct {
	Name string
	Age  *int
}



func BenchmarkSplite(b *testing.B) {
	b.StopTimer()
	var str = "#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//sqlArgRegex.FindAllString(str,-1)
		//strings.Split(str,"#{")
		//strings.SplitAfter()
		FindAllExpressConvertString(str)
	}
}

func TestFindAllExpressConvertString(t *testing.T) {
	var str = "#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}#{name}"
	fmt.Println(FindAllExpressConvertString(str))
}
