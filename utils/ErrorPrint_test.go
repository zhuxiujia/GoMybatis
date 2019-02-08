package utils

import (
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	var err = NewError("TestNewError", "aaa")
	fmt.Println(err)
}
