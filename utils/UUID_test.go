package utils

import (
	"fmt"
	"testing"
)

func TestCreateUUID(t *testing.T) {
	var id = CreateUUID()
	fmt.Println(id)
	if id == "" {
		t.Fatal("CreateUUID fail")
	}
}
