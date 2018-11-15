package GoMybatis

import "testing"

func TestManager(t *testing.T)  {
	var id ="1233214"
	var manager = DefaultTransationManager{}.New()
   var transcation,err=	manager.GetTransaction(nil,id)
   if err!=nil{
   	  panic(err)
   }
   transcation.
}