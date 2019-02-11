package GoMybatis

type ElementType = string

const (
	//root elements
	Element_ResultMap ElementType = "resultMap"
	Element_Insert    ElementType = "insert"
	Element_Delete    ElementType = "delete"
	Element_Update    ElementType = `update`
	Element_Select    ElementType = "select"
	Element_Sql       ElementType = "sql"

	//root templete elements
	Element_Insert_Templete ElementType = "insertTemplete"
	Element_Delete_Templete ElementType = "deleteTemplete"
	Element_Update_Templete ElementType = `updateTemplete`
	Element_Select_Templete ElementType = "selectTemplete"

	//child elements
	Element_bind      ElementType = "bind"
	Element_String    ElementType = "string"
	Element_If        ElementType = `if`
	Element_Trim      ElementType = "trim"
	Element_Foreach   ElementType = "foreach"
	Element_Set       ElementType = "set"
	Element_choose    ElementType = "choose"
	Element_when      ElementType = "when"
	Element_otherwise ElementType = "otherwise"
	Element_where     ElementType = "where"
	Element_Include   ElementType = "include"
)
