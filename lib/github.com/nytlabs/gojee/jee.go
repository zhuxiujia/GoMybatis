package jee

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	ZERO = iota
	CONST
	OP
	FUNC
	KEY
	K_START
	K_END
	Q_START
	Q_END
	SPACE
	NEXT
	D_STR
	S_STR
	ESC
	RESERVED
	EQ
)

var Ident = map[rune]int{
	'$':  FUNC,
	'.':  KEY,
	'+':  OP,
	'-':  OP,
	'/':  OP,
	'*':  OP,
	'!':  OP,
	'=':  OP,
	'>':  OP,
	'<':  OP,
	'&':  OP,
	'|':  OP,
	'(':  Q_START,
	')':  Q_END,
	'[':  K_START,
	']':  K_END,
	'"':  D_STR,
	'\'': S_STR,
	'\\': ESC,
	',':  NEXT,
}

var IdentStr = map[int]string{
	FUNC:     "FUNC",
	OP:       "OP",
	KEY:      "KEY",
	CONST:    "CONST",
	Q_START:  "Q_START",
	Q_END:    "Q_END",
	K_START:  "K_START",
	K_END:    "K_END",
	NEXT:     "NEXT",
	D_STR:    "D_STR",
	S_STR:    "S_STR",
	RESERVED: "RES",
	EQ:       "EQ",
}

type BMsg interface{}

type Token struct {
	Type  int
	Value string
}

type TokenTree struct {
	Type   int
	Value  interface{}
	Tokens []*TokenTree
	Parent *TokenTree
}

var tokenPopMap = map[int]func(rune, string) bool{
	D_STR: func(r rune, c string) bool {
		switch Ident[r] {
		case D_STR:
			return true
		}
		return false
	},
	S_STR: func(r rune, c string) bool {
		switch Ident[r] {
		case S_STR:
			return true
		}
		return false
	},
	KEY: func(r rune, c string) bool {
		switch getIdent(r) {
		case Q_START, Q_END, K_START, K_END, OP, FUNC, NEXT, KEY, D_STR, S_STR:
			return true
		}
		return false
	},
	OP: func(r rune, c string) bool {
		switch c {
		case "*", "+", "-", "/":
			return true
		}

		if len(c) >= 2 {
			return true
		}
		switch getIdent(r) {
		case Q_START, Q_END, K_START, K_END, FUNC, CONST, KEY, D_STR, S_STR, RESERVED:
			return true
		}
		return false
	},
	CONST: func(r rune, c string) bool {
		switch getIdent(r) {
		case Q_START, Q_END, K_START, K_END, OP, FUNC, NEXT, D_STR, S_STR, RESERVED:
			return true
		}
		return false
	},
	FUNC: func(r rune, c string) bool {
		switch getIdent(r) {
		case Q_START, Q_END, K_START, K_END, OP, FUNC, NEXT, KEY, D_STR, S_STR:
			return true
		}
		return false
	},
	RESERVED: func(r rune, c string) bool {
		switch getIdent(r) {
		case Q_START, Q_END, K_START, K_END, OP, FUNC, NEXT, KEY, D_STR, S_STR:
			return true
		}
		return false
	},
}

func getIdent(r rune) int {
	i, ok := Ident[r]
	switch {
	case ok:
		return i
	case unicode.IsNumber(r):
		return CONST
	case unicode.IsLetter(r) || unicode.IsPunct(r) || unicode.IsSymbol(r):
		return RESERVED
	case unicode.IsSpace(r):
		return SPACE
	}
	return ZERO
}

func emitToken(t []*Token, state int, value string) ([]*Token, string) {
	t = append(t, &Token{
		Type:  state,
		Value: value,
	})
	return t, ""
}

// probably should just use bufio.Scanner...
func Lexer(input string) ([]*Token, error) {
	var tokens []*Token
	var currWord string
	var state int
	var poppedStr bool
	var escaped bool

	for _, r := range input {

		// if we have a space and we aren't in a string
		if getIdent(r) == SPACE && state != D_STR && state != S_STR {
			continue
		}

		// if we have an escape char and we are in a string
		if getIdent(r) == ESC && (state == D_STR || state == S_STR) {
			escaped = true
			continue
		}

		if getIdent(r) == ZERO && state != D_STR && state != S_STR {
			return nil, errors.New(fmt.Sprintf("unexpected token: %s", string(r)))
		}

		switch state {
		case OP, FUNC, CONST, KEY, RESERVED:
			if tokenPopMap[state](r, currWord) {
				tokens, currWord = emitToken(tokens, state, currWord)
			}
		case D_STR, S_STR:
			if escaped {
				escaped = false
				break
			}
			if tokenPopMap[state](r, currWord) {
				currWord += string(r)
				tokens, currWord = emitToken(tokens, state, currWord)
				poppedStr = true
			} else {
				poppedStr = false
			}
		case Q_START, Q_END, K_START, K_END, NEXT:
			tokens, currWord = emitToken(tokens, state, currWord)
		}

		if !poppedStr {
			if len(currWord) == 0 {
				currWord = string(r)
				state = getIdent(r)
			} else {
				currWord += string(r)
			}
		} else {
			poppedStr = false
			state = ZERO
		}
	}

	if len(currWord) > 0 {
		tokens, _ = emitToken(tokens, state, currWord)
	}

	return tokens, nil
}

func buildTree(tokens []*Token) (*TokenTree, error) {
	var state int
	tree := &TokenTree{}
	var inKey bool // TODO: this needs to go
	var nested int
	var knested int

	// first pass:
	// take care of quantities, funcs, keys.
	for _, t := range tokens {

		item := &TokenTree{
			Value:  t.Value,
			Type:   t.Type,
			Parent: tree,
		}

		// this item should probably be in Lexer
		// convert value to float64 if number
		if item.Type == CONST {
			f, err := strconv.ParseFloat(t.Value, 64)
			if err != nil {

			} else {
				item.Value = f
			}
		}

		// this item should probably be in Lexer
		// get rid of quotes around our strings
		if item.Type == D_STR || item.Type == S_STR {
			item.Value = item.Value.(string)[1 : len(item.Value.(string))-1]
		}

		// this item should probably be in Lexer
		// create bool type
		if item.Type == RESERVED {
			switch item.Value {
			case "true", "false":
				f, err := strconv.ParseBool(t.Value)
				if err != nil {

				} else {
					item.Value = f
				}
			case "null":
				item.Value = nil
			default:
				return nil, errors.New(fmt.Sprintf("unexpected token: %s", item.Value))
			}
		}

		//
		if item.Type == K_START {
			item.Value = nil
		}

		// remove '.' from key name
		if item.Type == KEY {
			item.Value = item.Value.(string)[1:]
		}

		switch t.Type {
		case FUNC, CONST, RESERVED, D_STR, S_STR, NEXT:
			if inKey {
				for tree.Parent != nil && tree.Type == KEY && (tree.Type == KEY || tree.Type == K_START) {
					tree = tree.Parent
				}
				inKey = false
			}
			tree.Tokens = append(tree.Tokens, item)
		case OP:
			if inKey {
				for tree.Parent != nil && tree.Type == KEY && (tree.Type == KEY || tree.Type == K_START) {
					tree = tree.Parent
				}
				inKey = false
			}
			tree.Tokens = append(tree.Tokens, item)
		case KEY:
			tree.Tokens = append(tree.Tokens, item)
			if !inKey {
				tree = item
			}
			inKey = true
		case Q_START:
			nested++
			if state == FUNC {
				// we are a function
				tree = tree.Tokens[len(tree.Tokens)-1]
			} else {
				// we are a quantity
				tree.Tokens = append(tree.Tokens, item)
				tree = item
			}
		case K_START:
			if tree.Type != KEY && tree.Type != K_START {
				return nil, errors.New("unexpected [")
			}

			knested++
			tree.Tokens = append(tree.Tokens, item)
			tree = item

		case K_END:
			knested--
			// ???????
			if tree.Parent != nil {
				tree = tree.Parent
			} else {
				return nil, errors.New("unbalanced () or []")
			}
			inKey = true

		case Q_END:
			nested--
			if inKey {
				for tree.Parent != nil && (tree.Type == KEY || tree.Type == K_START) {
					tree = tree.Parent
				}
				inKey = false
			}

			if tree.Parent != nil {
				tree = tree.Parent
			} else {
				return nil, errors.New("unbalanced () or []")
			}
		}

		switch state {
		default:
		}

		state = t.Type
	}

	if nested != 0 || knested != 0 {
		return nil, errors.New("unbalanced () or []")
	}

	for tree.Parent != nil {
		tree = tree.Parent
	}

	return tree, nil
}

func not(tree *TokenTree) *TokenTree {
	var negate *TokenTree
	var newTokens []*TokenTree

	for _, t := range tree.Tokens {

		if t.Type == OP && t.Value == "!" {
			negate = t
			newTokens = append(newTokens, t)
			continue
		}

		if negate != nil {
			negate.Tokens = append(negate.Tokens, t)
			negate = nil
		} else {
			t.Parent = negate
			newTokens = append(newTokens, t)
		}

		if len(t.Tokens) > 0 {
			t = not(t)
		}
	}

	tree.Tokens = newTokens
	return tree
}

func split(tree *TokenTree, TokenType int, Values []string) *TokenTree {
	var nextTokens []*TokenTree

	popTokens := make([]*TokenTree, len(tree.Tokens))

	if len(tree.Tokens) > 0 {
		for i, t := range tree.Tokens {
			tree.Tokens[i] = split(t, TokenType, Values)
		}
	}

	copy(popTokens, tree.Tokens)

	for len(popTokens) > 2 {
		prev := popTokens[0]
		curr := popTokens[1]
		next := popTokens[2]

		if curr.Type == TokenType && inStringSlice(Values, curr.Value.(string)) {
			prev.Parent = curr
			next.Parent = curr
			curr.Tokens = append(curr.Tokens, prev, next)
			popTokens = popTokens[2:]
			popTokens[0] = curr
		} else {
			nextTokens = append(nextTokens, popTokens[0])
			popTokens = popTokens[1:]
		}
	}

	nextTokens = append(nextTokens, popTokens...)

	tree.Tokens = nextTokens

	return tree
}

func inStringSlice(a []string, b string) bool {
	for _, s := range a {
		if s == b {
			return true
		}
	}
	return false
}

func negative(tree *TokenTree) *TokenTree {
	var negate *TokenTree
	var newTokens []*TokenTree
	var state int

	for _, t := range tree.Tokens {

		// not entirely sure this case is correct
		if t.Type == OP && t.Value == "-" && state != CONST && state != KEY && state != Q_START {
			negate = t
			newTokens = append(newTokens, t)
			continue
		}

		if negate != nil {
			negate.Tokens = append(negate.Tokens, t)
			negate = nil
		} else {
			t.Parent = negate
			newTokens = append(newTokens, t)
		}

		state = t.Type

		if len(t.Tokens) > 0 {
			t = negative(t)
		}
	}

	tree.Tokens = newTokens
	return tree
}

func Parser(tokens []*Token) (*TokenTree, error) {

	tree, err := buildTree(tokens)
	if err != nil {
		return nil, err
	}

	tree = not(tree)
	tree = negative(tree)

	tree = split(tree, OP, []string{"&&", "||"})
	tree = split(tree, OP, []string{"*", "/"})
	tree = split(tree, OP, []string{"+", "-"})
	tree = split(tree, OP, []string{"==", ">=", ">", "<", "<=", "!="})

	return tree, nil
}

var opFuncsFloat = map[string]func(float64, float64) interface{}{
	"+": func(a float64, b float64) interface{} {
		return a + b
	},
	"-": func(a float64, b float64) interface{} {
		return a - b
	},
	"*": func(a float64, b float64) interface{} {
		return a * b
	},
	"/": func(a float64, b float64) interface{} {
		return a / b
	},
	"==": func(a float64, b float64) interface{} {
		return a == b
	},
	">=": func(a float64, b float64) interface{} {
		return a >= b
	},
	">": func(a float64, b float64) interface{} {
		return a > b
	},
	"<": func(a float64, b float64) interface{} {
		return a < b
	},
	"<=": func(a float64, b float64) interface{} {
		return a <= b
	},
	"!=": func(a float64, b float64) interface{} {
		return a != b
	},
}

var opFuncsString = map[string]func(string, string) interface{}{
	"+": func(a string, b string) interface{} {
		return a + b
	},
	"==": func(a string, b string) interface{} {
		return a == b
	},
	"!=": func(a string, b string) interface{} {
		return a != b
	},
}

var opFuncsBool = map[string]func(bool, bool) interface{}{
	"&&": func(a bool, b bool) interface{} {
		return a && b
	},
	"||": func(a bool, b bool) interface{} {
		return a || b
	},
	"==": func(a bool, b bool) interface{} {
		return a == b
	},
	"!=": func(a bool, b bool) interface{} {
		return a != b
	},
}

// this is a catch for types not bool, array, string, float
var opFuncsNil = map[string]func(interface{}, interface{}) interface{}{
	"==": func(a interface{}, b interface{}) interface{} {
		if a == nil && b == nil {
			return true
		}

		// comparing objects is a horrible condition and should be avoided
		return reflect.DeepEqual(a, b)
	},
	"!=": func(a interface{}, b interface{}) interface{} {
		return a != b
	},
}

var nullaryFuncs = map[string]func() (interface{}, error){
	"$now": func() (interface{}, error) {
		return float64(time.Now().UnixNano() / 1000 / 1000), nil
	},
}

var unaryFuncs = map[string]func(interface{}) (interface{}, error){
	"$sum": func(val interface{}) (interface{}, error) {
		valsArray, ok := val.([]interface{})
		if !ok {
			return nil, nil
		}
		sum := 0.0
		for _, i := range valsArray {
			sum += i.(float64)
		}
		return sum, nil
	},
	"$min": func(val interface{}) (interface{}, error) {
		valsArray, ok := val.([]interface{})
		if !ok {
			return nil, nil
		}

		min := valsArray[0].(float64)
		for i := 1; i < len(valsArray); i++ {
			min = math.Min(min, valsArray[i].(float64))
		}
		return min, nil
	},
	"$max": func(val interface{}) (interface{}, error) {
		valsArray, ok := val.([]interface{})
		if !ok {
			return nil, nil
		}

		max := valsArray[0].(float64)
		for i := 1; i < len(valsArray); i++ {
			max = math.Max(max, valsArray[i].(float64))
		}
		return max, nil
	},
	"$len": func(val interface{}) (interface{}, error) {
		valsArray, ok := val.([]interface{})
		if !ok {
			return nil, nil
		}

		return float64(len(valsArray)), nil
	},
	"$sqrt": func(val interface{}) (interface{}, error) {
		f, ok := val.(float64)
		if !ok || f < 0 {
			return nil, nil
		}

		return math.Sqrt(f), nil
	},
	"$abs": func(val interface{}) (interface{}, error) {
		f, ok := val.(float64)
		if !ok {
			return nil, nil
		}
		return math.Abs(f), nil
	},
	"$floor": func(val interface{}) (interface{}, error) {
		f, ok := val.(float64)
		if !ok {
			return nil, nil
		}

		return math.Floor(f), nil
	},
	"$keys": func(val interface{}) (interface{}, error) {
		var keyList []interface{}
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, nil
		}

		for k, _ := range m {
			keyList = append(keyList, k)
		}

		return keyList, nil
	},
	"$str": func(val interface{}) (interface{}, error) {
		switch v := val.(type) {
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64), nil
		case bool:
			if v {
				return "true", nil
			}
			return "false", nil
		case string:
			return v, nil
		case nil:
			return "null", nil
		case map[string]interface{}, []interface{}:
			b, err := json.Marshal(v)
			return string(b), err
		}
		return "", nil
	},
	"$num": func(val interface{}) (interface{}, error) {
		switch v := val.(type) {
		case float64:
			return v, nil
		case string:
			return strconv.ParseFloat(v, 64)
		case bool:
			if v {
				return 1, nil
			}
		}
		return 0.0, nil
	},
	"$~bool": func(val interface{}) (interface{}, error) {
		switch v := val.(type) {
		case []interface{}:
			if len(v) > 0 {
				return true, nil
			}
		case map[string]interface{}:
			return true, nil
		case float64:
			if math.IsNaN(v) {
				return false, nil
			}

			if v > 0 {
				return true, nil
			}
		case string:
			if len(v) > 0 {
				return true, nil
			}
		case bool:
			return v, nil
		}
		return false, nil
	},
	"$bool": func(val interface{}) (interface{}, error) {
		switch v := val.(type) {
		case string:
			return strconv.ParseBool(v)
		case bool:
			return v, nil
		}

		return nil, nil
	},
}

var binaryFuncs = map[string]func(interface{}, interface{}) (interface{}, error){
	"$parseTime": func(a interface{}, b interface{}) (interface{}, error) {
		layout, ok := a.(string)
		if !ok {
			return nil, nil
		}
		value, ok := b.(string)
		if !ok {
			return nil, nil
		}
		t, err := time.Parse(layout, value)
		if err != nil {
			return nil, err
		}
		return float64(t.UnixNano() / 1000 / 1000), nil
	},
	"$fmtTime": func(a interface{}, b interface{}) (interface{}, error) {
		layout, ok := a.(string)
		if !ok {
			return nil, nil
		}

		t, ok := b.(float64)
		if !ok {
			return nil, nil
		}

		return time.Unix(0, int64(time.Duration(t)*time.Millisecond)).Format(layout), nil
	},
	"$pow": func(a interface{}, b interface{}) (interface{}, error) {
		fa, ok := a.(float64)
		if !ok {
			return nil, nil
		}
		fb, ok := b.(float64)
		if !ok {
			return nil, nil
		}

		return math.Pow(fa, fb), nil
	},
	"$exists": func(a interface{}, b interface{}) (interface{}, error) {
		sb, ok := b.(string)
		if !ok {
			return nil, nil
		}

		ma, ok := a.(map[string]interface{})
		if !ok {
			return nil, nil
		}

		_, ok = ma[sb]
		if ok {
			return true, nil
		}
		return false, nil
	},
	"$contains": func(a interface{}, b interface{}) (interface{}, error) {
		sa, ok := a.(string)
		if !ok {
			return nil, nil
		}

		sb, ok := b.(string)
		if !ok {
			return nil, nil
		}
		return strings.Contains(sa, sb), nil
	},
	"$regex": func(a interface{}, b interface{}) (interface{}, error) {
		sa, ok := a.(string)
		if !ok {
			return nil, nil
		}

		sb, ok := b.(string)
		if !ok {
			return nil, nil
		}

		return regexp.MatchString(sb, sa)
	},
	"$has": func(a interface{}, b interface{}) (interface{}, error) {
		s, ok := a.([]interface{})
		if !ok {
			return nil, nil
		}

		for _, e := range s {
			switch c := e.(type) {
			case string:
				if c == b.(string) {
					return true, nil
				}
			case float64:
				if c == b.(float64) {
					return true, nil
				}
			case bool:
				if c == b.(bool) {
					return true, nil
				}
			default:
				if a == nil && b == nil {
					return true, nil
				}
			}
		}
		return false, nil
	},
}

func getKeyValues(t *TokenTree, input BMsg) (interface{}, error) {
	s, ok := t.Value.(string)

	if ok && len(s) > 0 {
		inputMap, ok := input.(map[string]interface{})
		if !ok {
			return nil, errors.New("could not assert to map")
		}

		input = inputMap[s]
	}

	var output []interface{}
	output = append(output, input)
	var accessed bool // this needs to be figured out!

	for _, sub := range t.Tokens {
		switch sub.Type {
		case K_START:
			switch c := sub.Value.(type) {
			case string:
				for j, _ := range output {
					outputMap, ok := output[j].(map[string]interface{})
					if !ok {
						return nil, errors.New("could not assert to map")
					}

					output[j] = outputMap[c]
				}
			case float64:
				for j, _ := range output {
					outputSlice, ok := output[j].([]interface{})
					if !ok {
						return nil, errors.New("could not assert to slice")
					}
					sliceIndex := int(c)
					if c < 0 || sliceIndex >= len(outputSlice) {
						output[j] = nil
					} else {
						output[j] = outputSlice[sliceIndex]
					}
				}
			default:
				accessed = true
				var newOutput []interface{}
				for j, _ := range output {
					arr, ok := output[j].([]interface{})
					if !ok {
						return nil, errors.New("could not assert to slice")
					}
					for _, e := range arr {
						newOutput = append(newOutput, e)
					}
				}
				output = newOutput
			}
		case KEY:
			for j, _ := range output {
				outputMap, ok := output[j].(map[string]interface{})
				if !ok {
					errors.New("could not assert to map")
				}

				subValue, ok := sub.Value.(string)
				if !ok {
					errors.New("invalid key for map")
				}
				output[j] = outputMap[subValue]
			}
		}
	}

	if len(output) == 1 && !accessed {
		return output[0], nil
	}

	return output, nil
}

func Eval(t *TokenTree, msg BMsg) (interface{}, error) {
	var tokenVal string

	switch t.Type {
	case OP, KEY, FUNC:
		_, ok := t.Value.(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("bad operation, key, or function: %s", t.Value))
		}
		tokenVal = t.Value.(string)
	}

	switch t.Type {
	case OP:
		if len(t.Tokens) == 1 {
			switch tokenVal {
			case "-":
				r, err := Eval(t.Tokens[0], msg)
				if err != nil {
					return nil, err
				}

				f, ok := r.(float64)
				if !ok {
					return nil, errors.New("cannot use - operator on non-number type")
				}

				return -1 * f, nil

			case "!":
				r, err := Eval(t.Tokens[0], msg)
				if err != nil {
					return nil, err
				}

				b, ok := r.(bool)
				if !ok {
					return nil, errors.New("cannot use ! operator on non-bool type")
				}

				return !b, nil
			}
			break
		}
		if len(t.Tokens) == 2 {
			a, err := Eval(t.Tokens[0], msg)
			if err != nil {
				return nil, err
			}

			b, err := Eval(t.Tokens[1], msg)
			if err != nil {
				return nil, err
			}
			// need to do comparisons for falsy-null || X
			// as well as != and ==

			switch ta := a.(type) {
			case float64:
				bf, ok := b.(float64)
				if !ok && tokenVal == "!=" {
					return true, nil
				} else if !ok && tokenVal == "==" {
					return false, nil
				} else if !ok {
					return nil, errors.New(fmt.Sprintf("cannot compare types: %s, %s", reflect.TypeOf(a), reflect.TypeOf(b)))
				}

				_, ok = opFuncsFloat[tokenVal]
				if !ok {
					return nil, errors.New(fmt.Sprintf("invalid operator for type: %s, %s", tokenVal, reflect.TypeOf(a)))
				}

				return opFuncsFloat[tokenVal](ta, bf), nil
			case string:
				bs, ok := b.(string)
				if !ok && tokenVal == "!=" {
					return true, nil
				} else if !ok && tokenVal == "==" {
					return false, nil
				} else if !ok {
					return nil, errors.New(fmt.Sprintf("cannot compare types: %s, %s", reflect.TypeOf(a), reflect.TypeOf(b)))
				}

				_, ok = opFuncsString[tokenVal]
				if !ok {
					return nil, errors.New(fmt.Sprintf("invalid operator for type: %s, %s", tokenVal, reflect.TypeOf(a)))
				}

				return opFuncsString[tokenVal](ta, bs), nil
			case bool:
				bb, ok := b.(bool)
				if !ok && tokenVal == "!=" {
					return true, nil
				} else if !ok && tokenVal == "==" {
					return false, nil
				} else if !ok {
					return nil, errors.New(fmt.Sprintf("cannot compare types: %s, %s", reflect.TypeOf(a), reflect.TypeOf(b)))
				}

				_, ok = opFuncsBool[tokenVal]
				if !ok {
					return nil, errors.New(fmt.Sprintf("invalid operator for type: %s, %s", tokenVal, reflect.TypeOf(a)))
				}

				return opFuncsBool[tokenVal](ta, bb), nil
			default:
				_, ok := opFuncsNil[tokenVal]
				if !ok {
					return nil, errors.New(fmt.Sprintf("invalid operator for type: %s, %s", tokenVal, reflect.TypeOf(a)))
				}

				return opFuncsNil[tokenVal](a, b), nil
			}
		}
	case S_STR, D_STR, CONST, RESERVED:
		return t.Value, nil
	case KEY:
		input := msg

		for _, sub := range t.Tokens {
			if len(sub.Tokens) > 0 {
				key, err := Eval(sub.Tokens[0], input)
				if err != nil {
					return nil, err
				}

				switch key.(type) {
				case string:
					sub.Type = KEY
				}
				sub.Value = key
				sub.Tokens = nil
			}
		}

		return getKeyValues(t, input)
	case FUNC:
		if len(t.Tokens) == 0 {
			_, ok := nullaryFuncs[tokenVal]
			if !ok {
				return nil, errors.New(fmt.Sprintf("func does not exist or wrong num of arguments: %s", tokenVal))
			}
			return nullaryFuncs[tokenVal]()
		}
		if len(t.Tokens) == 1 {
			a, err := Eval(t.Tokens[0], msg)
			if err != nil {
				return nil, err
			}

			_, ok := unaryFuncs[tokenVal]
			if !ok {
				return nil, errors.New(fmt.Sprintf("func does not exist or wrong num of arguments: %s", tokenVal))
			}

			return unaryFuncs[tokenVal](a)
		} else if len(t.Tokens) == 3 {

			a, err := Eval(t.Tokens[0], msg)
			if err != nil {
				return nil, err
			}

			b, err := Eval(t.Tokens[2], msg)
			if err != nil {
				return nil, err
			}

			_, ok := binaryFuncs[tokenVal]
			if !ok {
				return nil, errors.New(fmt.Sprintf("func does not exist or wrong num of arguments: %s", tokenVal))
			}

			return binaryFuncs[tokenVal](a, b)
		}
		return nil, errors.New(fmt.Sprintf("func does not exist or wrong num of arguments: %s", tokenVal))
	default:
		if len(t.Tokens) > 0 {
			return Eval(t.Tokens[0], msg)
		}
	}

	return nil, nil
}

func FmtTokens(tl []*Token) {
	for _, t := range tl {
		fmt.Printf("(" + IdentStr[t.Type] + " " + t.Value + ") ")
	}
}

func FmtTokenTree(tree *TokenTree, d int) {
	fmt.Printf("\n")
	for i := 0; i < d; i++ {
		fmt.Printf("  ")
	}

	fmt.Printf("[")
	if tree.Type != ZERO {
		fmt.Printf("%s ", IdentStr[tree.Type])
	}
	fmt.Printf("%s", tree.Value)
	d++
	for _, t := range tree.Tokens {
		FmtTokenTree(t, d)
	}
	fmt.Printf("]")
}
