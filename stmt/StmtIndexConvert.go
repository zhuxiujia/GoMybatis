package stmt

// stmt convert
// example mysql: input 1 -> ?
// oracle : input 0 ->   :val1
// sqlite: input 0 ->   " ? "
type StmtIndexConvert interface {
	Convert(index int) string
}
