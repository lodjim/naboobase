package utils

import (
	"fmt"
	"reflect"
)

func Get(fieldName string, s interface{}) (interface{}, error) {
	val := reflect.ValueOf(s)

	// If s is a pointer, dereference it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure s is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct or a pointer to a struct, got %v", val.Kind())
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field '%s' does not exist", fieldName)
	}

	return field.Interface(), nil
}

func Set(fieldName string, value interface{}, s interface{}) error {
	val := reflect.ValueOf(s)

	// If s is a pointer, dereference it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		return fmt.Errorf("expected a pointer to a struct, got %v", val.Kind())
	}

	// Ensure s is a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct or a pointer to a struct, got %v", val.Kind())
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field '%s' does not exist", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", fieldName)
	}

	fieldValue := reflect.ValueOf(value)
	if field.Type() != fieldValue.Type() {
		return fmt.Errorf("value type for field '%s' is %s, got %s", fieldName, field.Type(), fieldValue.Type())
	}

	field.Set(fieldValue)
	return nil
}

func GetTaggedFields(s interface{}, tag string) []string {
	var uniqueFields []string
	val := reflect.ValueOf(s)

	// If s is a pointer, dereference it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure s is a struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct")
		return nil
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue, ok := field.Tag.Lookup("db")
		if ok && tagValue == tag {
			uniqueFields = append(uniqueFields, field.Name)
		}
	}

	return uniqueFields
}
