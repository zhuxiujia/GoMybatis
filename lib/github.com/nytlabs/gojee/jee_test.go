package jee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

type Test struct {
	exp    string
	result string
}

var Tests = []Test{
	Test{
		exp:    `.int == 5`,
		result: `true`,
	},
	Test{
		exp:    `.float == 5.5`,
		result: `true`,
	},
	Test{
		exp:    `.["es'c\"ape'.key"]`,
		result: `null`,
	},
	Test{
		exp:    `.['escape.key']`,
		result: `{"nested":{"foo.bar":"baz"}}`,
	},
	Test{
		exp:    `.['escape.key']['nested']`,
		result: `{"foo.bar":"baz"}`,
	},
	Test{
		exp:    `.['escape.key']['nested']['foo.bar']`,
		result: `"baz"`,
	},
	Test{
		exp:    `.['escape.key'].nested["foo.bar"]`,
		result: `"baz"`,
	},
	Test{
		exp:    `.string == "hello world"`,
		result: `true`,
	},
	Test{
		exp:    `.int -- .int`,
		result: `10`,
	},
	Test{
		exp:    `.float - -.float`,
		result: `11`,
	},
	Test{
		exp:    `.int +-.float`,
		result: `-0.5`,
	},
	Test{
		exp:    `.int/-.float`,
		result: `-0.9090909090909091`,
	},
	Test{
		exp:    `-.float/-.float`,
		result: `1`,
	},
	Test{
		exp:    `-.float/.float`,
		result: `-1`,
	},
	Test{
		exp:    `.bool == false`,
		result: `true`,
	},
	Test{
		exp:    `.nil == null`,
		result: `true`,
	},
	Test{
		exp:    `null == .nil`,
		result: `true`,
	},
	Test{
		exp:    `.nested.foo.zip`,
		result: `"zap"`,
	},
	Test{
		exp:    `.arrayInt`,
		result: `[1,2,3,4,5,6,7,8,9,10]`,
	},
	Test{
		exp:    `.arrayFloat`,
		result: `[1.1,2.2,3.3,4.4,5.5,6.6,7.7,8.8,9.9,10]`,
	},
	Test{
		exp:    `.arrayInt[0]`,
		result: `1`,
	},
	Test{
		exp:    `.arrayObj[0].nested[0].id`,
		result: `"foo"`,
	},
	Test{
		exp:    `.arrayInt[]`,
		result: `[1,2,3,4,5,6,7,8,9,10]`,
	},
	Test{
		exp:    `.arrayObj[]`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.arrayObj`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.arrayObj[].name`,
		result: `["foo","bar","baz"]`,
	},
	Test{
		exp:    `.arrayObj[].val`,
		result: `[2,2.5,10]`,
	},
	Test{
		exp:    `.arrayObj[].array[]`,
		result: `[1,2,3,1,2,3,7,8,9]`,
	},
	Test{
		exp:    `.arrayObj[].array`,
		result: `[[1,2,3],[1,2,3],[7,8,9]]`,
	},
	Test{
		exp:    `.arrayObj[0].array`,
		result: `[1,2,3]`,
	},
	Test{
		exp:    `(true && false) && true == false`,
		result: `true`,
	},
	Test{
		exp:    `$sum(.arrayInt)`,
		result: `55`,
	},
	Test{
		exp:    `$sum(.arrayInt) + $sum(.arrayInt[])`,
		result: `110`,
	},
	Test{
		exp:    `$sum(.arrayFloat) < 100.0`,
		result: `true`,
	},
	Test{
		exp:    `$sum(.arrayObj[].array[]) == 36`,
		result: `true`,
	},
	Test{
		exp:    `true && true && false || (!true) == !true`,
		result: `true`,
	},
	Test{
		exp:    `true && false && (true && (true && false))`,
		result: `false`,
	},
	Test{
		exp:    `!(.int == 2 + 3) == false`,
		result: `true`,
	},
	Test{
		exp:    `100 - ((3/2)*20 + 7 -8)`,
		result: `71`,
	},
	Test{
		exp:    `     100.0 -          ( (   3/2 )*20 + 7 -8  )`,
		result: `71`,
	},
	Test{
		exp:    `"hello" + " " + "world"`,
		result: `"hello world"`,
	},
	Test{
		exp:    `true && true && (true || false)`,
		result: `true`,
	},
	Test{
		exp:    `(true || false) && true && true`,
		result: `true`,
	},
	Test{
		exp:    `false`,
		result: `false`,
	},
	Test{
		exp:    `true`,
		result: `true`,
	},
	Test{
		exp:    `null`,
		result: `null`,
	},
	Test{
		exp:    `.['escape.key']`,
		result: `{"nested":{"foo.bar":"baz"}}`,
	},
	Test{
		exp:    `.arrayInt[0]`,
		result: `1`,
	},
	Test{
		exp:    `.['arrayInt'][0]`,
		result: `1`,
	},
	Test{
		exp:    `.arrayObj[0].nested`,
		result: `[{"id":"foo","no":"zoo"}]`,
	},
	Test{
		exp:    `.['arrayObj'][0]['nested']`,
		result: `[{"id":"foo","no":"zoo"}]`,
	},
	Test{
		exp:    `.arrayObj[]`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.['arrayObj'][]['nested'][]['id']`,
		result: `["foo","zof","zif"]`,
	},
	Test{
		exp:    `.arrayObj[].array[2]`,
		result: `[3,3,9]`,
	},
	Test{
		exp:    `.arrayObj[1].nested[].id`,
		result: `["zof"]`,
	},
	Test{
		exp:    `.arrayFloat`,
		result: `[1.1,2.2,3.3,4.4,5.5,6.6,7.7,8.8,9.9,10]`,
	},
	Test{
		exp:    `.arrayFloat[]`,
		result: `[1.1,2.2,3.3,4.4,5.5,6.6,7.7,8.8,9.9,10]`,
	},
	Test{
		exp:    `.arrayObj`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.arrayObj[]`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.arrayObj[].array`,
		result: `[[1,2,3],[1,2,3],[7,8,9]]`,
	},
	Test{
		exp:    `.arrayObj[].array[]`,
		result: `[1,2,3,1,2,3,7,8,9]`,
	},
	Test{
		exp:    `.arrayObj[].array[1]`,
		result: `[2,2,8]`,
	},
	Test{
		exp:    `.arrayObj[].nested[].id`,
		result: `["foo","zof","zif"]`,
	},
	Test{
		exp:    `.arrayObj[0].nested[0].id`,
		result: `"foo"`,
	},
	Test{
		exp:    `.arrayObj[2].nested[].id`,
		result: `["zif"]`,
	},
	Test{
		exp:    `.['arrayFloat']`,
		result: `[1.1,2.2,3.3,4.4,5.5,6.6,7.7,8.8,9.9,10]`,
	},
	Test{
		exp:    `.['arrayFloat'][]`,
		result: `[1.1,2.2,3.3,4.4,5.5,6.6,7.7,8.8,9.9,10]`,
	},
	Test{
		exp:    `.['arrayObj']`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.['arrayObj'][]`,
		result: `[{"array":[1,2,3],"bool":false,"hasKey":true,"name":"foo","nested":[{"id":"foo","no":"zoo"}],"nil":null,"sameNum":10,"sameStr":"all","val":2},{"array":[1,2,3],"bool":true,"name":"bar","nested":[{"id":"zof","no":"fum"}],"nil":null,"sameNum":10,"sameStr":"all","val":2.5},{"array":[7,8,9],"bool":false,"name":"baz","nested":[{"id":"zif","no":"zaf"}],"nil":null,"sameNum":10,"sameStr":"all","val":10}]`,
	},
	Test{
		exp:    `.['arrayObj'][]['array']`,
		result: `[[1,2,3],[1,2,3],[7,8,9]]`,
	},
	Test{
		exp:    `.['arrayObj'][]['array'][]`,
		result: `[1,2,3,1,2,3,7,8,9]`,
	},
	Test{
		exp:    `.['arrayObj'][]['array'][1]`,
		result: `[2,2,8]`,
	},
	Test{
		exp:    `.['arrayObj'][]['nested']`,
		result: `[[{"id":"foo","no":"zoo"}],[{"id":"zof","no":"fum"}],[{"id":"zif","no":"zaf"}]]`,
	},
	Test{
		exp:    `.['arrayObj'][]['nested'][]`,
		result: `[{"id":"foo","no":"zoo"},{"id":"zof","no":"fum"},{"id":"zif","no":"zaf"}]`,
	},
	Test{
		exp:    `.['arrayObj'][]['nested'][]['id']`,
		result: `["foo","zof","zif"]`,
	},
	Test{
		exp:    `.['arrayObj'][0]['nested'][0]['id']`,
		result: `"foo"`,
	},
	Test{
		exp:    `.['arrayObj'][2]['nested'][]['id']`,
		result: `["zif"]`,
	},
	Test{
		exp:    `$len(.arrayObj)`,
		result: `3`,
	},
	Test{
		exp:    `$len(.arrayObj[])`,
		result: `3`,
	},
	Test{
		exp:    `$len(.['arrayObj'][]['array'])`,
		result: `3`,
	},
	Test{
		exp:    `$len(.['arrayObj'][]['array'][])`,
		result: `9`,
	},
	Test{
		exp:    `-(2 * (-2 * (-5 + (-5))))`,
		result: `-40`,
	},
	Test{
		exp:    `$pow(.int,2)`,
		result: `25`,
	},
	Test{
		exp:    `$sqrt(100)`,
		result: `10`,
	},
	Test{
		exp:    `$sqrt($sum(.arrayInt))`,
		result: `7.416198487095663`,
	},
	Test{
		exp:    `$pow( (-0.1) * 10, 2)`,
		result: `1`,
	},
	Test{
		exp:    `$abs(-100)`,
		result: `100`,
	},
	Test{
		exp:    `$abs(100)`,
		result: `100`,
	},
	Test{
		exp:    `$max(.arrayInt)`,
		result: `10`,
	},
	Test{
		exp:    `$min(.arrayInt)`,
		result: `1`,
	},
	Test{
		exp:    `$min(.arrayObj[].array[])`,
		result: `1`,
	},
	Test{
		exp:    `$floor(.arrayFloat[0])`,
		result: `1`,
	},
	Test{
		exp:    `$contains(.string,"hello")`,
		result: `true`,
	},
	Test{
		exp:    `$contains("http://en.wikipedia.org/wiki/List_of_animals_with_fraudulent_diplomas","wikipedia")`,
		result: `true`,
	},
	Test{
		exp:    `$contains("http://en.wikipedia.org/wiki/List_of_animals_with_fraudulent_diplomas","dogs")`,
		result: `false`,
	},
	// needs a special test ... $keys returns an unordered list
	//Test{
	//            exp:`$keys(.arrayObj[0])`,
	//            result:`["val","name","sameStr","hasKey","nil","array","nested","bool","sameNum"]`,
	//},
	Test{
		exp:    `$has($keys(.), "arrayString")`,
		result: `true`,
	},
	Test{
		exp:    `$has($keys(.), "nope")`,
		result: `false`,
	},
	Test{
		exp:    `$exists(., "arrayString")`,
		result: `true`,
	},
	Test{
		exp:    `$exists(., "nope")`,
		result: `false`,
	},
	Test{
		exp:    `$has(.arrayFloat, 1.1)`,
		result: `true`,
	},
	Test{
		exp:    `$has($keys(.), "arrayString") || $has($keys(.), "nope") `,
		result: `true`,
	},
	Test{
		exp:    `$has($keys(.), "arrayString") && $has($keys(.), "nope") `,
		result: `false`,
	},
	Test{
		exp:    `.#_k__`,
		result: `1`,
	},
	Test{
		exp:    `$num(.float_str) == 5.123131`,
		result: `true`,
	},
	Test{
		exp:    `$str($num(.float_str)) == .float_str`,
		result: `true`,
	},
	Test{
		exp:    `$parseTime("Mon Jan 2 15:04:05 -0700 MST 2006","Wed Jan 1 00:00:00 +0000 GMT 2014") == 1388534400000`,
		result: `true`,
	},
	Test{
		exp:    `$num($fmtTime("2006", $parseTime("Mon Jan 2 15:04:05 -0700 MST 2006","Wed Jan 1 00:00:00 +0000 GMT 2014"))) == 2014`,
		result: `true`,
	},
	Test{
		exp:    `$now() > 1388534400000`,
		result: `true`,
	},
	Test{
		exp:    `$num($fmtTime("2006", $now())) > 2006`,
		result: `true`,
	},
	Test{
		exp:    `$str(.float_str) == "5.123131"`,
		result: `true`,
	},
	Test{
		exp:    `$num(.int) == 5`,
		result: `true`,
	},
	Test{
		exp:    `$num(.bool) == 0`,
		result: `true`,
	},
	Test{
		exp:    `$num(.a) == 0`,
		result: `true`,
	},
	Test{
		exp:    `$num(.empty) == 0`,
		result: `true`,
	},
	Test{
		exp:    `$bool("true") && true`,
		result: `true`,
	},
	Test{
		exp:    `$bool("false") && true`,
		result: `false`,
	},
	Test{
		exp:    `$bool(1)`,
		result: `null`,
	},
	Test{
		exp:    `$bool(null)`,
		result: `null`,
	},
	Test{
		exp:    `$~bool(null)`,
		result: `false`,
	},
	Test{
		exp:    `$~bool(.empty)`,
		result: `false`,
	},
	Test{
		exp:    `$~bool(.a.b.c)`,
		result: `true`,
	},
	Test{
		exp:    `$~bool("asdsajdasd")`,
		result: `true`,
	},
	Test{
		exp:    `$~bool(1)`,
		result: `true`,
	},
}

func TestAll(t *testing.T) {
	var umsg BMsg

	testFile, _ := ioutil.ReadFile("test.json")

	json.Unmarshal(testFile, &umsg)

	for _, test := range Tests {
		tokenized, err := Lexer(test.exp)
		if err != nil {
			t.Error("failed lex")
		}

		tree, err := Parser(tokenized)
		if err != nil {
			t.Error("failed parse")
		}

		result, err := Eval(tree, umsg)

		if err != nil {
			t.Error("failed eval")
		}

		var rmsg BMsg
		err = json.Unmarshal([]byte(test.result), &rmsg)
		if err != nil {
			t.Error(err, "bad test")
		}

		if reflect.DeepEqual(rmsg, result) {
			fmt.Println("\x1b[32;1mOK\x1b[0m", test.exp)
		} else {
			t.Fail()
			fmt.Println("\x1b[31;1m", "FAIL", "\x1b[0m", rmsg, "\t", result)
			fmt.Println("Expected Value", rmsg, "\tResult Value:", result)
			fmt.Println("Expected Type: ", reflect.TypeOf(rmsg), "\tResult Type:", reflect.TypeOf(result))
		}
	}
}

func BenchmarkJSON(b *testing.B) {
	var umsg BMsg
	testFile, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(testFile, &umsg)
	tokenized, _ := Lexer(`.['arrayObj'][2]['nested'][]['id']`)
	tree, _ := Parser(tokenized)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(tree, umsg)
	}
}

func BenchmarkMath(b *testing.B) {
	var umsg BMsg
	testFile, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(testFile, &umsg)
	tokenized, _ := Lexer(`100 * -($sum(.arrayInt) + 5)`)
	tree, _ := Parser(tokenized)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(tree, umsg)
	}
}

func BenchmarkRegex(b *testing.B) {
	var umsg BMsg
	testFile, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(testFile, &umsg)
	tokenized, _ := Lexer(`$regex(.string, "hello*")`)
	tree, _ := Parser(tokenized)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(tree, umsg)
	}
}

func BenchmarkContains(b *testing.B) {
	var umsg BMsg
	testFile, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(testFile, &umsg)
	tokenized, _ := Lexer(`$contains(.string, "hello")`)
	tree, _ := Parser(tokenized)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(tree, umsg)
	}
}
