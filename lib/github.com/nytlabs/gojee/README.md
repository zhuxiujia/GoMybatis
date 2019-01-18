# jee 
**jee** (json expression evaluator) transforms JSON  through logical and mathematical expressions. jee can be used from the command line or as a Go package. It is inspired by the fantastic (and much more fully featured) [jq]("http://stedolan.github.io/jq/"). 

jee was created out of the need for a simple JSON query language in [streamtools]("https://github.com/nytlabs/streamtools/"). jee is designed for stream processing and provides a reusable token tree. 

####get the library

    go get github.com/nytlabs/gojee

####make and install the binary

    cd $gopath/src/github.com/nytlabs/gojee/jee
    go install


### usage (binary)
##### querying JSON

get the entire input object:

    > echo '{"a": 3}' | jee '.'
    {
        "a": 3
    }

get a value for a specific key:

    > echo '{"a": 3, "b": 4}' | jee '.a'
    3

get a value from an array:

    > echo '{"a": [4,5,6]}' | jee '.a[0]'
    4

get all values from an array:

    > echo '{"a": [4,5,6]}' | jee '.a[]'
    [
        4,
        5,
        6
    ]

query all objects inside array 'a' for key 'id':

    > echo '{"a": [{"id":"foo"},{"id":"bar"},{"id":"baz"}]}' | jee '.a[].id'
    [
        "foo",
        "bar",
        "baz"
    ]


##### arithmetic 
\+ - * /

    > echo '{"a": 10}' | jee '(.a * 100)/-10 * 5'
    -500
    
##### comparison 
\> >= < <= !=

    > echo '{"a": 10}' | jee '(.a * 100)/-10 * 5 == -500'
    true
    > echo '{"a": 10}' | jee '(.a * 100)/-10 * 5 > 0'
    false

##### logical
|| &&
    
    > echo '{"a": false}' | jee '!(.a && true) || false  == true'
    true
    
##### functions

###### types

**`$num(x {bool, float64, string, nil})`**
<br />
Converts `x` to a float64. If `x` is a bool, 1 is returned for true and 0 for false. If `x` is nil, 0 is returned. 
<br /><br />
**`$str(x {bool, float64, string, nil, object, []*))`**
<br />
Converts `x` to a string. If `x` is a bool, "true" is returned for true and "false" for false. "null" is returned for nil. If `x` is an object or an array, it is marshaled into a JSON string. 
<br /><br />
**`$bool(x {bool, string})`**
<br />
Converts `x` to a bool. See [strconv.ParseBool](http://golang.org/pkg/strconv/#ParseBool)
<br /><br />
**`$~bool(x {bool, float64, string, nil, object, []*})`**
<br />
Truthy conversion of `x` to a bool. Falsey values:`null`,`NaN`,`0`,`false`, and `arrays with a length of 0`. 
<br /><br />
###### math

**`$sqrt(x float64)`**
<br />
Returns square root of `x`.
<br /><br />
**`$pow(x float64, y float64)`**
<br />
Returns `x`^`y`.
<br /><br />
**`$floor(x float64)`**
<br />
Returns nearest downward integer for `x`.
<br /><br />
**`$abs(x float64)`**
<br />
Returns absolute value of `x`.
<br /><br />
###### arrays

**`$len(a []interface{})`**
<br />
Returns the length of array `a`. 
<br /><br />
**`$has( a {[]bool, []float64, []string, []nil}, val {bool, float64, string, nil} )`**
<br />
Checks to see if array `a` contains `val`. Returns bool. `val` cannot be an object.
<br /><br />
**`$sum(a []float64)`**
<br />
Returns the sum of array `a`.
<br /><br />
**`$min(a []float64)`**
<br />
Returns the minumum of array `a`.
<br /><br />
**`$max(a []float64)`**
<br />
Returns the maximum of array `a`.
<br /><br />
###### objects

**`$keys(o object)`**
<br />
Returns an array of keys in object `o`.
<br /><br />
**`$exists(o object, key string)`**
<br />
Checks to see if `key` exists in map `o`. Returns bool. `$exists()` does a map lookup and is faster than `$has($keys(o), "foo")`
<br /><br /><br />
###### date and time

**`$now()`**
<br />
Returns current system time in float64 (epoch milliseconds).
<br /><br />
**`$parseTime(layout string, t string)`**
<br />
Accepts a time layout in golang [time format](http://golang.org/pkg/time/#pkg-constants). t is parsed and returned as epoch milliseconds in float64.
<br /><br />
**`$fmtTime(layout string, t float64)`**
<br />
Accepts a time layout in golang [time format](http://golang.org/pkg/time/#pkg-constants). t is expected in epoch milliseconds. Returns a formatted string. 
<br /><br />
###### strings

**`$contains(s string, substr string)`**
<br />
see [strings.Contains](http://golang.org/pkg/strings/#Contains)
<br /><br />
**`$regex(pattern string, s string)`**
<br />
see [regexp.MatchString](http://golang.org/pkg/regexp/#MatchString). Much slower than `$contains()`
<br /><br />
see `jee_test.go` for examples.

### package usage
#####`Lexer(string) []*Token, error`
converts a jee query string into a slice of tokens

#####`Parser([]*Tokens) *TokenTree, error`
builds a parse tree out token slice from `Lexer()`

#####`Eval(*TokenTree, {}interface) {}interface, error`
evaluates a variable of type interface{} with a *TokenTree generated from `Parser()`. Only types given by [`json.Unmarshal`]("http://golang.org/pkg/encoding/json/#Unmarshal") are supported.

### quirks
* Types are strictly enforced. `false || "foo"` will produce a type error.
* `null` and `0` are not falsey
* Using a JSON key as an array index or an escaped key in bracket notation will not currently be evaluated. ie: `.a[.b]`
* All numbers in a jee query must start with a digit. numbers <1 should start with a 0. use `0.1` instead of `.1`
* Bracket notation is available for keys that need escaping `.["foo"]["bar"]`]
* Queries for JSON keys or indices that do not exist return `null` (to test if a key exists, use `$exists`)
* jee does not support variables, conditional expressions, or assignment 
* jee may be very quirky in general.

### changes
- **.0.1.1** addition of $bool, $~bool, $num, $str, $now, $fmtTime, $parseTime. Fix for non-alphanumeric characters in JSON keys. 
- **.0.1.0** initial release