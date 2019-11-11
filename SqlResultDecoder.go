package GoMybatis

//sql查询结果解码
type SqlResultDecoder interface {
	//resultMap = in xml resultMap element
	//dbData = select the SqlResult
	//decodeResultPtr = need decode result type
	Decode(resultMap map[string]*ResultProperty, SqlResult QueryResult, decodeResultPtr interface{}) error
}
