package logar

type Map map[string]any

type TypeKind string

const (
	TypeKind_Text     TypeKind = "text"
	TypeKind_Int      TypeKind = "int"
	TypeKind_Float    TypeKind = "float"
	TypeKind_Bool     TypeKind = "bool"
	TypeKind_Time     TypeKind = "time"
	TypeKind_Duration TypeKind = "duration"
)

const (
	LogarLogs Model = "logar.logs"
)
