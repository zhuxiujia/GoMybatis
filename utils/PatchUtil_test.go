package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestFixTestExpressionSymbol(t *testing.T) {
	var bytes = []byte(`<if test="page <= 0 and size != 0">limit #{page}, #{size}</if>`)
	FixTestExpressionSymbol(&bytes)
	fmt.Println(string(bytes))
	if strings.Index(string(bytes), "&lt;") == -1 {
		t.Fatal("TestFixTestExpressionSymbol fail")
	}
}
