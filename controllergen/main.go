package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/lodjim/naboobase/utils"
)

const controllerTemplate = `package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lodjim/naboobase/core"
	"github.com/lodjim/naboobase/models"
)

// Create{{.Model}} creates a new {{.Model}} in the database
func Create{{.Model}}(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateCreateHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.{{.Model}}Request{} },
		NewModel:    func() interface{} { return &models.{{.Model}}{} },
		NewResponse: func() interface{} { return &models.{{.Model}}Response{} },
		Functionality: "{{.Collection}}",
		Collection:  "{{.Collection}}",
		Preprocess:  nil,
	})
}

func GetOne{{.Model}}(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetOneHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.{{.Model}}Request{} },
		NewModel:    func() interface{} { return &models.{{.Model}}{} },
		NewResponse: func() interface{} { return &models.{{.Model}}Response{} },
		Functionality: "{{.Collection}}",
		Collection:  "{{.Collection}}",
		Preprocess:  nil,
	})
}
func GetAll{{.Model}}(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetAllHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.{{.Model}}Request{} },
		NewModel:    func() interface{} { return &models.{{.Model}}{} },
		NewResponse: func() interface{} { return &models.{{.Model}}Response{} },
		Functionality: "{{.Collection}}",
		Collection:  "{{.Collection}}",
		Preprocess:  nil,
	})
}

func Update{{.Model}}(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateUpdateHandler(db, core.HandlerConfig{
			NewRequest:  func() interface{} { return &models.{{.Model}}Request{} },
			NewModel:    func() interface{} { return &models.{{.Model}}{} },
			NewResponse: func() interface{} { return &models.{{.Model}}Response{} },
			Functionality: "{{.Collection}}",
			Collection:  "{{.Collection}}",
			Preprocess:  nil,
	})
}

func Delete{{.Model}}(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateDeleteHandler(db, core.HandlerConfig{
			NewRequest:  func() interface{} { return &models.{{.Model}}Request{} },
			NewModel:    func() interface{} { return &models.{{.Model}}{} },
			NewResponse: func() interface{} { return &models.{{.Model}}Response{} },
			Functionality: "{{.Collection}}",
			Collection:  "{{.Collection}}",
			Preprocess:  nil,
	})
}


func init() {
	core.AutoEndpointFuncRegistry["{{.Collection}}-POST"] = Create{{.Model}}
	core.AutoEndpointFuncRegistry["{{.Collection}}-GET-ID"] = GetOne{{.Model}}
	core.AutoEndpointFuncRegistry["{{.Collection}}-GET"] = GetAll{{.Model}}
	core.AutoEndpointFuncRegistry["{{.Collection}}-PUT-ID"] = Update{{.Model}}
	core.AutoEndpointFuncRegistry["{{.Collection}}-DELETE-ID"] = Delete{{.Model}}
}
`

type ControllerConfig struct {
	Model      string
	Collection string
}

func generateControllerFile(config ControllerConfig) error {
	// Create filename
	filename := fmt.Sprintf("./controllers/%s_controller.go", strings.ToLower(config.Model))

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create template
	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Execute template
	err = tmpl.Execute(file, config)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("Generated controller file: %s\n", filename)
	return nil
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "generate" {
		fmt.Println("Usage: <executable> generate")
		os.Exit(1)
	}
	logger := log.New(os.Stdout, "PROTOC_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	jsonDir, err := os.ReadDir("models")
	if err != nil {
		logger.Fatalf("Failed to read 'json' directory: %s", err.Error())
	}
	for _, filename := range jsonDir {
		base := strings.Split(filename.Name(), ".")[0]
		if strings.Contains(base, "request") || strings.Contains(base, "response") || strings.Contains(base, "user") {
			continue
		}
		config := ControllerConfig{
			Model:      utils.ConvertToCamelCase(base),
			Collection: base,
		}
		err := generateControllerFile(config)
		if err != nil {
			fmt.Printf("Error generating controller: %v\n", err)
		}
	}
}
