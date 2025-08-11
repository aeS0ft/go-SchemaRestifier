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
	TypeStringVar   types = "varchar(255)"
	TypeText        types = "text"
	TypeVarchar     types = "varchar"
	TypeChar        types = "char"
	TypeSerial      types = "serial"
	TypeBigSerial   types = "bigserial"
	TypeBigInt      types = "bigint"
	TypeSmallInt    types = "smallint"
	TypeNumeric     types = "numeric"
	TypeUUID        types = "uuid"
	TypeJSONB       types = "jsonb"
	TypeBytea       types = "bytea"
	TypeTimestamp   types = "timestamp"
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
	TypeStringVar:   "string",
	TypeText:        "string",
	TypeVarchar:     "string",
	TypeChar:        "string",
	TypeSerial:      "int",
	TypeBigSerial:   "int64",
	TypeBigInt:      "int64",
	TypeSmallInt:    "int16",
	TypeNumeric:     "float64",
	TypeUUID:        "string",
	TypeJSONB:       "map[string]interface{}",
	TypeBytea:       "[]byte",
	TypeTimestamp:   "time.Time",
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
	case TypeString, TypeInteger, TypeBoolean, TypeFloat, TypeObject, TypeArray, TypeDate, TypeDateTime, TypeChar, TypeSerial, TypeBigSerial, TypeBigInt, TypeSmallInt, TypeNumeric, TypeUUID, TypeJSONB, TypeBytea, TypeTimestamp:
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
	case "varchar(255)":
		return TypeStringVar, true
	case "text":
		return TypeText, true
	case "varchar":
		return TypeVarchar, true
	case "char":
		return TypeChar, true
	case "serial":
		return TypeSerial, true
	case "bigserial":
		return TypeBigSerial, true
	case "timestamp":
		return TypeTimestamp, true
	default:
		return "", false
	}
}
