package utils

import (
	"fmt"
	"naboobase/models" // Update this import path to match your module name
	"reflect"
	"strings"
)

// GetModelType takes a model name (like "user") and returns the corresponding type (like User)
func GetModelType(modelName string) (reflect.Type, error) {
	if len(modelName) == 0 {
		return nil, fmt.Errorf("empty model name")
	}

	// Convert first letter to uppercase to match Go type naming convention
	typeName := strings.ToUpper(modelName[:1]) + strings.ToLower(modelName[1:])

	// Create an instance of a models struct to get the correct package info
	dummy := &models.User{}
	pkg := reflect.TypeOf(dummy).Elem().PkgPath()

	// Look for the type in the models package
	pkgPath := fmt.Sprintf("%s.%s", pkg, typeName)

	// Get the actual type from the models package
	modelType := reflect.TypeOf(models.User{})
	if modelType.Name() == typeName {
		return modelType, nil
	}

	return nil, fmt.Errorf("model %s not found: type %s not found", typeName, pkgPath)
}
