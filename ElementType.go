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

	//root template elements
	Element_Insert_Template ElementType = "insertTemplate"
	Element_Delete_Template ElementType = "deleteTemplate"
	Element_Update_Template ElementType = `updateTemplate`
	Element_Select_Template ElementType = "selectTemplate"

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

func isMethodElement(tag ElementType) bool {
	switch tag {
	case Element_Insert, Element_Delete, Element_Update, Element_Select,
		Element_Insert_Template, Element_Delete_Template, Element_Update_Template, Element_Select_Template:
		return true
	}
	return false
}
