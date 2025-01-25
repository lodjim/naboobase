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
	DBTag      string // New field for db tag
	Validation string
}

func ConvertToCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i := range words {
		if len(words[i]) > 0 {
			words[i] = strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	return strings.Join(words, "")
}

func ConvertToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) {
			// Add an underscore before the uppercase letter (except for the first character)
			if i > 0 {
				result = append(result, '_')
			}
			// Convert the uppercase letter to lowercase
			result = append(result, unicode.ToLower(r))
		} else {
			// Append the lowercase letter as is
			result = append(result, r)
		}
	}
	return string(result)
}

func ParseStruct(name string, data map[string]interface{}, structs map[string]*StructDefinition) *StructDefinition {
	st := &StructDefinition{
		Name:   name,
		Fields: make([]FieldDefinition, 0),
	}
	structs[name] = st

	for key, value := range data {
		field := FieldDefinition{
			Name:    ToGoFieldName(key),
			JSONTag: key,
			BSONTag: key,
			DBTag:   "", // Initialize DBTag
		}

		switch v := value.(type) {
		case map[string]interface{}:
			if val, ok := v["value"]; ok {
				// Handle fields with validation
				field.Type = GetGoType(val)
				if validationRule, ok := v["validation"].(string); ok {
					field.Validation = validationRule
				} else {
					// Only set default validation if no explicit validation is provided
					field.Validation = GetDefaultValidation(field.Type, key)
				}
				if dbTag, ok := v["db"].(string); ok {
					field.DBTag = dbTag
					// Override type if db tag is "autogenerate"
					if dbTag == "autogenerate" {
						field.Type = "primitive.ObjectID"
					}
				}
			} else {

				nestedName := fmt.Sprintf("%s%s", name, field.Name)
				ParseStruct(nestedName, v, structs)
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

func GetDefaultValidation(fieldType, fieldName string) string {
	var validations []string

	// Add validation rules based on field type or name
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

func ToGoFieldName(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		switch strings.ToLower(p) {
		case "url":
			parts[i] = "URL"
		case "id":
			parts[i] = "ID"
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

func GenerateFile(structs map[string]*StructDefinition, packageName string) []byte {
	var buf bytes.Buffer

	// Start with a package declaration
	buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Generate each struct
	for _, st := range structs {
		GenerateStructCode(&buf, st)
	}

	return buf.Bytes()
}
