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

// ParseTypes converts a string to the corresponding types enum value.
func ParseTypes(s string) (types, bool) {
	switch s {
	case "string":
		return TypeString, true
	case "integer":
		return TypeInteger, true
	case "bool":
		return TypeBoolean, true
	case "float64":
		return TypeFloat, true
	case "map[string]interface{}":
		return TypeObject, true
	case "[]interface{}":
		return TypeArray, true
	case "date":
		return TypeDate, true
	case "datetime":
		return TypeDateTime, true
	case "[integer]":
		return TypeIntArray, true
	case "[string]":
		return TypeStringArray, true
	default:
		return "", false
	}
}
