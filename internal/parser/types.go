package parser

import (
	"go-SchemaRestifier/internal/datastructures"
	"strings"
)

// Package parser provides functionality to parse and handle schemas in a structured format.
type CrudOperation struct {
	Operation    string `json:"operation"`
	Endpoint     string `json:"endpoint"`
	AuthRequired bool   `json:"auth_required"`
	By           string `json:"by,omitempty"`
}

// Schema represents the structure of a parsed schema.
type Schema struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Columns     *[]Column       `json:"columns"` // dynamic keys/values
	Crud        []CrudOperation `json:"crud"`    // dynamic keys/values
}

// Column represents a single column in a schema.
type Column struct {
	Name          string               `json:"name"`                  // Name of the column
	Type          string               `json:"type"`                  // Type of the column (e.g., string, integer, etc.)
	Description   string               `json:"description,omitempty"` // Optional description of the column
	PrimaryKey    bool                 `json:"primary_key,omitempty"` // Indicates if the column is a primary key
	Nullable      bool                 `json:"nullable,omitempty"`
	Unique        bool                 `json:"unique,omitempty"`
	Hidden        bool                 `json:"hidden,omitempty"`
	Nestedcolumns *datastructures.Node `json:"json_data,omitempty"` // Optional field for nested schemas
	Capabilities  map[string]bool      `json:"query,omitempty"`     // e.g., {"select": true, "filter": false}
}

// Types represents supported type strings.
type Types string

const (
	TypeString      Types = "string"
	TypeInteger     Types = "integer"
	TypeBoolean     Types = "boolean"
	TypeFloat       Types = "float"
	TypeObject      Types = "object"
	TypeArray       Types = "array"
	TypeDate        Types = "date"
	TypeDateTime    Types = "datetime"
	TypeIntArray    Types = "[integer]"
	TypeStringArray Types = "[string]"
	TypeStringVar   Types = "varchar(255)"
	TypeText        Types = "text"
	TypeVarchar     Types = "varchar"
	TypeChar        Types = "char"
	TypeSerial      Types = "serial"
	TypeBigSerial   Types = "bigserial"
	TypeBigInt      Types = "bigint"
	TypeSmallInt    Types = "smallint"
	TypeNumeric     Types = "numeric"
	TypeUUID        Types = "uuid"
	TypeJSONB       Types = "jsonb"
	TypeBytea       Types = "bytea"
	TypeTimestamp   Types = "timestamp"
)

var typeNames = map[Types]string{
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
func (t Types) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	}
	return "unknown"
}

// ftypes returns the canonical type or "unknown" if not found.
func ftypes(s Types) Types {
	switch s {
	case TypeString, TypeInteger, TypeBoolean, TypeFloat, TypeObject, TypeArray, TypeDate, TypeDateTime, TypeChar, TypeSerial, TypeBigSerial, TypeBigInt, TypeSmallInt, TypeNumeric, TypeUUID, TypeJSONB, TypeBytea, TypeTimestamp:
		return s
	default:
		return "unknown"
	}
}

// ParseTypes converts a string to the corresponding types enum value.
func ParseTypes(s string) (Types, bool) {
	switch {
	case s == "string":
		return TypeString, true
	case s == "integer":
		return TypeInteger, true
	case s == "bool":
		return TypeBoolean, true
	case s == "float64":
		return TypeFloat, true
	case s == "map[string]interface{}":
		return TypeObject, true
	case s == "[]interface{}":
		return TypeArray, true
	case s == "date":
		return TypeDate, true
	case s == "datetime":
		return TypeDateTime, true
	case s == "[integer]":
		return TypeIntArray, true
	case s == "[string]":
		return TypeStringArray, true
	case strings.Contains(s, "varchar"):
		return TypeStringVar, true
	case s == "text":
		return TypeText, true
	case s == "varchar":
		return TypeVarchar, true
	case s == "char":
		return TypeChar, true
	case s == "serial":
		return TypeSerial, true
	case s == "bigserial":
		return TypeBigSerial, true
	case s == "timestamp":
		return TypeTimestamp, true
	}
	return "unknown", false
}
