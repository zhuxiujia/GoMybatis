package GoExpress

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	Parser("a = 'b'")
}

func TestParser2(t *testing.T) {
	var opts = ParserOperators("a +1 > b")
	fmt.Println(opts)
}
