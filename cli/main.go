package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lodjim/naboobase/utils"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "generate" {
		fmt.Println("Usage: <executable> generate")
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "PROTOC_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

	jsonDir, err := os.ReadDir("json")
	if err != nil {
		logger.Fatalf("Failed to read 'json' directory: %s", err.Error())
	}

	for _, filename := range jsonDir {
		jsonPath := fmt.Sprintf("./json/%s", filename.Name())
		logger.Printf("Processing file: %s", jsonPath)
		inputFile := jsonPath
		// Use the file name (without extension) to derive the main struct name
		base := strings.Split(filename.Name(), ".")[0]
		outputFile := fmt.Sprintf("./models/%s.go", base)
		packageName := "models"
		structName := utils.ConvertToCamelCase(base)

		// Read and parse JSON
		jsonData, err := ioutil.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(jsonData, &data); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}

		// Parse structs (the main struct and any nested structs)
		structs := make(map[string]*utils.StructDefinition)
		utils.ParseStruct(structName, data, structs)

		// Check if MongoDB primitive import is needed
		needsPrimitive := false
		for _, st := range structs {
			for _, field := range st.Fields {
				if field.Type == "primitive.ObjectID" {
					needsPrimitive = true
					break
				}
			}
			if needsPrimitive {
				break
			}
		}

		// Generate code with proper imports and struct definitions
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))
		if needsPrimitive {
			buf.WriteString("import \"go.mongodb.org/mongo-driver/bson/primitive\"\n\n")
		}

		// Write all struct definitions
		for _, st := range structs {
			utils.GenerateStructCode(&buf, st)
		}

		// Append an init function that registers the main struct.
		// This init function will be executed automatically when the package is loaded.
		// Format and write the generated code
		formattedCode, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Printf("Error formatting code: %v\n", err)
			os.Exit(1)
		}

		if err := ioutil.WriteFile(outputFile, formattedCode, 0644); err != nil {
			fmt.Printf("Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully generated struct definitions for", structName)
	}
	logger.Println("Processing completed.")
}
