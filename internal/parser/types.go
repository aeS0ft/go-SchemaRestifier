package parser

// Package parser provides functionality to parse and handle schemas in a structured format.

// Schema represents the structure of a parsed schema.
type Schema struct {
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	Fields        *[]Column                `json:"columns"` // dynamic keys/values
	Crud          map[string]interface{}   `json:"crud"`
	Nestedcolumns map[string][]interface{} `json:"nested"` // Optional field for nested schemas
}

// Column represents a single column in a schema.
type Column struct {
	Name        string `json:"name"`                  // Name of the column
	Type        string `json:"type"`                  // Type of the column (e.g., string, integer, etc.)
	Description string `json:"description,omitempty"` // Optional description of the column
	PrimaryKey  bool   `json:"primary_key,omitempty"` // Indicates if the column is a primary key

}
