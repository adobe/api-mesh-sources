package main

import (
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

const (
	failed int = iota
	success
)

func main() {
	args := os.Args[3:]
	schema := os.Args[1]
	prefix := os.Args[2]

	if len(args) != 0 {
		for _, v := range args {
			schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schema))
			documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s%s", prefix, v))
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				fmt.Println("Schema Validation Error", err)
			}

			if !result.Valid() {
				fmt.Printf("The connecot %s is not valid. see errors :\n", v)
				for _, desc := range result.Errors() {
					fmt.Printf("- %s\n", desc)
				}
				os.Exit(failed)
			}
		}
	}
	fmt.Println("Validation success")
	os.Exit(success)
}
