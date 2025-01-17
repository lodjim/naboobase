package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"naboobase/utils"
	"os"
	"strings"
)

func main() {
	if os.Args[1] == "generate" {
		logger := log.New(os.Stdout, "PROTOC_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

		jsonDir, err := os.ReadDir("json")
		if err != nil {
			logger.Fatalf("Failed to read 'proto' directory: %s", err.Error())
		}

		for _, filename := range jsonDir {
			jsonPath := fmt.Sprintf("./json/%s", filename.Name())
			logger.Printf("Processing file: %s", jsonPath)
			inputFile := jsonPath
			outputFile := fmt.Sprintf("./models/%s.go", strings.Split(filename.Name(), ".")[0])
			packageName := "models"
			structName := utils.ConvertToCamelCase(strings.Split(filename.Name(), ".")[0])

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

			// Generate code
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

			structs := make(map[string]*utils.StructDefinition)
			utils.ParseStruct(structName, data, structs)

			// Write all struct definitions
			for _, st := range structs {
				utils.GenerateStructCode(&buf, st)
			}

			// Format and write the code
			formattedCode, err := format.Source(buf.Bytes())
			if err != nil {
				fmt.Printf("Error formatting code: %v\n", err)
				os.Exit(1)
			}

			if err := ioutil.WriteFile(outputFile, formattedCode, 0644); err != nil {
				fmt.Printf("Error writing output file: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Successfully generated struct definitions")

		}

		logger.Println("Processing completed.")
	}
}
