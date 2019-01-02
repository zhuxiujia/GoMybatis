package utils

import "testing"

func TestCreateUUID(t *testing.T) {
	var id = CreateUUID()
	if id == "" {
		t.Fatal("CreateUUID fail")
	}
}
