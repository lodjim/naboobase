package utils

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type StructDefinition struct {
	Name   string
	Fields []FieldDefinition
}

type FieldDefinition struct {
	Name       string
	Type       string
	JSONTag    string
	BSONTag    string
	DBTag      string
	Validation string
}

type EnumDefinition struct {
	Type   string   // The Go type of the enum (e.g., "string", "int")
	Values []string // The possible values as strings
}

// ConvertToCamelCase converts a snake_case string to CamelCase.
func ConvertToCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i := range words {
		if len(words[i]) > 0 {
			words[i] = strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	return strings.Join(words, "")
}

// ConvertToSnakeCase converts a CamelCase string to snake_case.
func ConvertToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// ParseStruct parses a JSON object into a struct definition.
func ParseStruct(name string, data map[string]interface{}, structs map[string]*StructDefinition, enums map[string]EnumDefinition) *StructDefinition {
	st := &StructDefinition{
		Name:   name,
		Fields: make([]FieldDefinition, 0),
	}
	structs[name] = st

	for key, value := range data {
		if key == "_config" {
			continue
		}
		field := FieldDefinition{
			Name:    ToGoFieldName(key),
			JSONTag: key,
			BSONTag: key,
			DBTag:   "", // Initialize DBTag
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Check if it's an enum
			if typ, ok := v["type"].(string); ok && typ == "enum" {
				if vals, ok := v["values"].([]interface{}); ok && len(vals) > 0 {
					// Determine the enum type based on the first value
					enumType := GetGoType(vals[0])
					enumName := fmt.Sprintf("%s%sEnum", name, field.Name)
					strVals := make([]string, len(vals))
					for i, val := range vals {
						strVals[i] = fmt.Sprintf("%v", val)
					}
					// Store the enum definition
					enums[enumName] = EnumDefinition{Type: enumType, Values: strVals}
					field.Type = enumName
					// Add validation
					field.Validation = "oneof=" + strings.Join(strVals, " ")
				}
			} else if val, ok := v["value"]; ok {
				// Handle fields with validation
				field.Type = GetGoType(val)
				if validationRule, ok := v["validation"].(string); ok {
					field.Validation = validationRule
				} else {
					field.Validation = GetDefaultValidation(field.Type, key)
				}
				if dbTag, ok := v["db"].(string); ok {
					field.DBTag = dbTag
					if dbTag == "autogenerate" {
						field.Type = "primitive.ObjectID"
					}
				}
			} else {
				// Nested struct
				nestedName := fmt.Sprintf("%s%s", name, field.Name)
				ParseStruct(nestedName, v, structs, enums)
				field.Type = nestedName
				field.Validation = "required"
			}
		default:
			field.Type = GetGoType(value)
			field.Validation = GetDefaultValidation(field.Type, key)
		}

		st.Fields = append(st.Fields, field)
	}

	return st
}

// GenerateStructCode generates the Go code for a struct.
func GenerateStructCode(buf *bytes.Buffer, st *StructDefinition) {
	buf.WriteString(fmt.Sprintf("type %s struct {\n", st.Name))
	for _, field := range st.Fields {
		tags := fmt.Sprintf("`json:\"%s\" bson:\"%s\"", field.JSONTag, field.BSONTag)
		if field.DBTag != "" {
			tags += fmt.Sprintf(" db:\"%s\"", field.DBTag)
		}
		if field.Validation != "" {
			tags += fmt.Sprintf(" validate:\"%s\"", field.Validation)
		}
		tags += "`"
		buf.WriteString(fmt.Sprintf("\t%s %s %s\n", field.Name, field.Type, tags))
	}
	buf.WriteString("}\n\n")
}

// GetGoType determines the Go type for a given value.
func GetGoType(v interface{}) string {
	switch val := v.(type) {
	case bool:
		return "bool"
	case float64:
		if val == float64(int64(val)) {
			return "int"
		}
		return "float64"
	case string:
		return "string"
	case []interface{}:
		if len(val) > 0 {
			return "[]" + GetGoType(val[0])
		}
		return "[]interface{}"
	case nil:
		return "interface{}"
	default:
		return "interface{}"
	}
}

// GetDefaultValidation provides default validation rules based on field type and name.
func GetDefaultValidation(fieldType, fieldName string) string {
	var validations []string

	switch {
	case fieldType == "string":
		if strings.Contains(strings.ToLower(fieldName), "email") {
			validations = append(validations, "email")
		} else if strings.Contains(strings.ToLower(fieldName), "url") {
			validations = append(validations, "url")
		} else {
			validations = append(validations, "max=255")
		}
	case fieldType == "int":
		validations = append(validations, "gte=0")
	}

	return strings.Join(validations, ",")
}

// ToGoFieldName converts a JSON key to a Go field name (e.g., snake_case to CamelCase).
func ToGoFieldName(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		switch strings.ToLower(p) {
		case "url":
			parts[i] = "URL"
		case "id":
			parts[i] = "Id"
		default:
			r := []rune(p)
			if len(r) > 0 {
				r[0] = unicode.ToUpper(r[0])
			}
			parts[i] = string(r)
		}
	}
	return strings.Join(parts, "")
}

// GenerateFile generates the complete Go file with enums and structs.
func GenerateFile(structs map[string]*StructDefinition, enums map[string]EnumDefinition, packageName string, needsPrimitive bool) []byte {
	var buf bytes.Buffer

	// Package declaration (must come first)
	buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Imports (if needed, comes after package declaration)
	if needsPrimitive {
		buf.WriteString("import \"go.mongodb.org/mongo-driver/bson/primitive\"\n\n")
	}

	// Generate enums
	for enumName, enumDef := range enums {
		buf.WriteString(fmt.Sprintf("type %s %s\n\n", enumName, enumDef.Type))
		buf.WriteString("const (\n")
		for _, val := range enumDef.Values {
			constName := fmt.Sprintf("%s%s", enumName, ToGoFieldName(val))
			if enumDef.Type == "string" {
				buf.WriteString(fmt.Sprintf("\t%s %s = \"%s\"\n", constName, enumName, val))
			} else {
				buf.WriteString(fmt.Sprintf("\t%s %s = %s\n", constName, enumName, val))
			}
		}
		buf.WriteString(")\n\n")
	}

	// Generate structs
	for _, st := range structs {
		GenerateStructCode(&buf, st)
	}

	return buf.Bytes()
}

// Example usage (not part of the package, for illustration)
func GenerateCode(data map[string]interface{}) []byte {
	structs := make(map[string]*StructDefinition)
	enums := make(map[string]EnumDefinition)
	ParseStruct("MyStruct", data, structs, enums)
	return GenerateFile(structs, enums, "mypackage", false)
}
