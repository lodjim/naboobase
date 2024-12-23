package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if os.Args[1] == "generate" {
		logger := log.New(os.Stdout, "PROTOC_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

		protoDir, err := os.ReadDir("json")
		if err != nil {
			logger.Fatalf("Failed to read 'proto' directory: %s", err.Error())
		}

		for _, filename := range protoDir {
			protoPath := fmt.Sprintf("./proto/%s", filename.Name())
			logger.Printf("Processing file: %s", protoPath)

			cmd := exec.Command("protoc", "--go_out=.", protoPath)
			res, err := cmd.Output()
			if err != nil {
				logger.Printf("Error running protoc on %s: %s", protoPath, err.Error())
				continue
			}

			logger.Printf("Successfully processed %s: %s", protoPath, string(res))
		}

		logger.Println("Processing completed.")
	}
}
