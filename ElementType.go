package GoMybatis

type ElementType = string

const (
	Element_String  = "string"
	Element_If      = `if`
	Element_Trim    = "trim"
	Element_Foreach = "foreach"
	Element_Set     = "set"
	Element_choose = "choose"
	Element_when = "when"
	Element_otherwise = "otherwise"
)
