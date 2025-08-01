package generator

// types represents supported type strings.
type types string

const (
	TypeString      types = "string"
	TypeInteger     types = "integer"
	TypeBoolean     types = "boolean"
	TypeFloat       types = "float"
	TypeObject      types = "object"
	TypeArray       types = "array"
	TypeDate        types = "date"
	TypeDateTime    types = "datetime"
	TypeIntArray    types = "[integer]"
	TypeStringArray types = "[string]"
)

var typeNames = map[types]string{
	TypeString:      "string",
	TypeInteger:     "int",
	TypeBoolean:     "bool",
	TypeFloat:       "float64",
	TypeObject:      "map[string]interface{}",
	TypeArray:       "[]interface{}",
	TypeDate:        "date",
	TypeDateTime:    "datetime",
	TypeIntArray:    "[]int",
	TypeStringArray: "[]string",
}

// String returns the string representation of the types value.
func (t types) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	}
	return "unknown"
}

// ftypes returns the canonical type or "unknown" if not found.
func ftypes(s types) types {
	switch s {
	case TypeString, TypeInteger, TypeBoolean, TypeFloat, TypeObject, TypeArray, TypeDate, TypeDateTime:
		return s
	default:
		return "unknown"
	}
}
