package GoMybatis

type RootElementType = string

const (
	Element_ResultMap RootElementType = "resultMap"
	Element_Insert    RootElementType = "insert"
	Element_Delete    RootElementType = "delete"
	Element_Update    RootElementType = `update`
	Element_Select    RootElementType = "select"

	Element_bind   RootElementType = "bind"
)
