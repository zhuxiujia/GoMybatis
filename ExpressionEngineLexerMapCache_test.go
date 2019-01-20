package GoMybatis

import (
	"fmt"
	"testing"
)

func TestExpressionEngineLexerMapCache_Get(t *testing.T) {
	engine := ExpressionEngineLexerMapCache{}.New()
	err := engine.Set("foo", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExpressionEngineLexerMapCache_Set(t *testing.T) {
	engine := ExpressionEngineLexerMapCache{}.New()
	err := engine.Set("foo", 1)
	if err != nil {
		t.Fatal(err)
	}
	result, err := engine.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
