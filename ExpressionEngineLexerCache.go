package GoMybatis

//Lexer 结果缓存
type ExpressionEngineLexerCache interface {
	Set(expression string, lexer interface{}) error
	Get(expression string) (interface{}, error)
}
