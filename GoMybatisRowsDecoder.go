package GoMybatis

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func rows2maps(rows *sql.Rows) ( string, error) {
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list := "["
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("log:", err)
			panic(err.Error())
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])

			cell := fmt.Sprintf(`"%v":"%v"`, columName, value)
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	list = list[0 : len(list)-1]
	list += "]"
	fmt.Println(list)
	return list,nil
}
//	fields, err := rows.Columns()
//	if err != nil {
//		return nil, err
//	}
//	for rows.Next() {
//		result, err := row2map(rows, fields)
//		if err != nil {
//			return nil, err
//		}
//		resultsSlice = append(resultsSlice, result)
//	}
//	return resultsSlice, nil
//}

func row2map(rows *sql.Rows, fields []string) (resultsMap map[string]interface{}, err error) {
	result := make(map[string]interface{})
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}
	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			result[key] = []byte{}
			continue
		}
		if data, err := value2Bytes(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}
	}
	return result, nil
}
func value2Bytes(rawValue *reflect.Value) (interface{}, error) {
	if rawValue.IsValid() && rawValue.CanInterface() {
		return rawValue.Interface(), nil
	}
	return nil, nil
}
