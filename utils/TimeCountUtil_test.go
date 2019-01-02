package utils

import (
	"testing"
	"time"
)

func TestDurationToString(t *testing.T) {
	var str = DurationToString(time.Second)
	if str != "s" {
		t.Fatal("TestDurationToString fail")
	}
}
