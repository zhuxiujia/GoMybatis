package GoMybatis

import (
	"database/sql"
	"reflect"
	"strconv"
	"fmt"
	"time"
	"github.com/zhuxiujia/GoMybatis/lib/github.com/satori/go.uuid"
)

type MysqlEngine struct {
	SessionEngine
	DB *sql.DB
}



func (this MysqlEngine) NewSession() *Session {
	uuids, _ := uuid.NewV4()
	var uuidstrig = uuids.String()
	var isCommitedOrRollbacked = false
	var mysqlLocalSession = LocalSqlSession{
		SessionId:              uuidstrig,
		db:                     this.DB,
		isCommitedOrRollbacked: &isCommitedOrRollbacked,
	}
	var session = Session(&mysqlLocalSession)
	return &session
}

//打开一个本地引擎
func Open(driverName, dataSourceName string) (*SessionEngine, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	var mysqlEngine = MysqlEngine{
		DB: db,
	}
	var engine = SessionEngine(mysqlEngine)
	return &engine, nil
}


//打开一个本地引擎
func OpenRemote(addr string,RetryTime int) (*SessionEngine, error) {
	var TransationRMClient = TransationRMClient{
		RetryTime: RetryTime,
		Addr:      addr,
	}
	var engine = RemoteSessionEngine{}.New(&TransationRMClient)
	var sessionEngine=SessionEngine(&engine)
	return &sessionEngine, nil
}

//bean 工厂，根据xml配置创建函数,并且动态代理到你定义的struct func里
//bean 参数必须为指针类型,指向你定义的struct
//你定义的struct必须有可导出的func属性,例如：
//type MyUserMapperImpl struct {
//	UserMapper                                                 `mapperPath:"/mapper/user/UserMapper.xml"`
//	SelectById    func(id string, result *model.User) error    `mapperParams:"id"`
//	SelectByPhone func(id string, phone string, result *model.User) error `mapperParams:"id,phone"`
//	DeleteById    func(id string, result *int64) error         `mapperParams:"id"`
//	Insert        func(arg model.User, result *int64) error
//}
//func的参数支持2种函数，第一种函数 基本参数个数无限制(并且需要用Tag指定参数名逗号隔开,例如`mapperParams:"id,phone"`)，最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//func的参数支持2种函数，第二种函数第一个参数必须为结构体(例如 arg model.User,该结构体的属性可以指定tag `json:"xxx"`为参数名称),最后一个参数必须为返回数据类型的指针(例如result *model.User)，返回值为error
//使用UseProxyMapper函数设置代理后即可正常使用。
func UseProxyMapperBySessionEngine(bean interface{}, xml []byte, engine *SessionEngine) {
	UseProxyMapperByEngine(bean, xml, engine)
}

//-------------------------------------------------------------工具

func rows2maps(rows *sql.Rows) (resultsSlice []map[string][]byte, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2map(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}
	return resultsSlice, nil
}

func row2map(rows *sql.Rows, fields []string) (resultsMap map[string][]byte, err error) {
	result := make(map[string][]byte)
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
func value2Bytes(rawValue *reflect.Value) ([]byte, error) {
	str, err := value2String(rawValue)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func value2String(rawValue *reflect.Value) (str string, err error) {
	aa := reflect.TypeOf((*rawValue).Interface())
	vv := reflect.ValueOf((*rawValue).Interface())
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.FormatInt(vv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = strconv.FormatUint(vv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
	case reflect.String:
		str = vv.String()
	case reflect.Array, reflect.Slice:
		switch aa.Elem().Kind() {
		case reflect.Uint8:
			data := rawValue.Interface().([]byte)
			str = string(data)
			if str == "\x00" {
				str = "0"
			}
		default:
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
		// time type
	case reflect.Struct:
		if aa.ConvertibleTo(TimeType) {
			str = vv.Convert(TimeType).Interface().(time.Time).Format(time.RFC3339Nano)
		} else {
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	case reflect.Bool:
		str = strconv.FormatBool(vv.Bool())
	case reflect.Complex128, reflect.Complex64:
		str = fmt.Sprintf("%v", vv.Complex())
		/* TODO: unsupported types below
		   case reflect.Map:
		   case reflect.Ptr:
		   case reflect.Uintptr:
		   case reflect.UnsafePointer:
		   case reflect.Chan, reflect.Func, reflect.Interface:
		*/
	default:
		err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
	}
	return
}
