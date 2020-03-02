package ast

import (
	"fmt"
	"testing"
)

func BenchmarkFindRawExpressString(b *testing.B) {
	b.StopTimer()
	var str = "${name1}${name2}${name3}${name4}"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		FindRawExpressString(str)
	}
}

func BenchmarkFindExpress(b *testing.B) {
	b.StopTimer()
	var str = "#{name1}#{name2}#{name3}#{name4}"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		FindExpress(str)
	}
}

func TestFindExpress(t *testing.T) {
	var str = "#{name1}#{name2}#{name3}#{name4}"
	for i := 0; i < 10000; i++ {
		var result = FindExpress(str)
		if !(result[0] == "name1" && result[1] == "name2" && result[2] == "name3" && result[3] == "name4") {
			panic("FindExpress fail not equal!")
		}
	}
}

func TestFindRawExpressString(t *testing.T) {
	var str = "${name1}${name2}${name3}${name4}"
	for i := 0; i < 10000; i++ {
		var result = FindRawExpressString(str)
		fmt.Println(result)
		if !(result[0] == "name1" && result[1] == "name2" && result[2] == "name3" && result[3] == "name4") {
			panic("FindExpress fail not equal!")
		}
	}
}
