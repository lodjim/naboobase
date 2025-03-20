package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lodjim/naboobase/models" // Update this import path to match your module name
)

func GetModelType(modelName string) (reflect.Type, error) {
	if len(modelName) == 0 {
		return nil, fmt.Errorf("empty model name")
	}
	typeName := strings.ToUpper(modelName[:1]) + strings.ToLower(modelName[1:])
	dummy := &models.User{}
	pkg := reflect.TypeOf(dummy).Elem().PkgPath()
	pkgPath := fmt.Sprintf("%s.%s", pkg, typeName)
	modelType := reflect.TypeOf(models.User{})
	if modelType.Name() == typeName {
		return modelType, nil
	}

	return nil, fmt.Errorf("model %s not found: type %s not found", typeName, pkgPath)
}
