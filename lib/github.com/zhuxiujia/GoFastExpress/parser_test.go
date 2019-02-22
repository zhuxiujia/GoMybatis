package GoFastExpress

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var e error
	//_,e=Parser("a = 'b'")
	//if e!=nil{
	//	t.Fatal(e)
	//}
	_, e = Parser(" a += b")
	if e != nil {
		t.Fatal(e)
	}
}

func TestParser2(t *testing.T) {
	var opts = ParserOperators("a +1 > b")
	fmt.Println(opts)
}

func BenchmarkParser(b *testing.B) {
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, e := Parser(" a + b")
		if e != nil {
			b.Fatal(e)
		}
	}
}
