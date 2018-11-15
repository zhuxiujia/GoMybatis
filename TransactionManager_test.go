package GoMybatis

import (
	"testing"
	"fmt"
)

func TestManager(t *testing.T)  {
	var id ="1233214"
	engine, err := Open("mysql", "*?charset=utf8&parseTime=True&loc=Local") //此处请按格式填写你的mysql链接，这里用*号代替
	if err != nil {
		panic(err.Error())
	}
	var SessionFactory = SessionFactory{}.New(engine)
	var fac=TransactionFactory{}.New()
	var manager = DefaultTransationManager{}.New(&SessionFactory,&fac)

	//start
    transcation,err :=	manager.GetTransaction(nil,id,id)
   if err!=nil{
   	  panic(err)
   }
   fmt.Println(transcation)
}