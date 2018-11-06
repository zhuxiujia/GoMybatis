package GoMybatis

import "regexp"

var reg = regexp.MustCompile("\\s+")
func ReplaceAllBlankSpace(sql string) string {
	return reg.ReplaceAllString(sql, " ")
}
