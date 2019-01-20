package GoMybatis

type ExpressionEngineLexerCacheable interface {
	SetUseLexerCache(use bool) error
	LexerCacheable() bool
}
