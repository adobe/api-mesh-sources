package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type ValidateConnector struct {
	rootPath       string
	schemaPath     string
	connectorPaths []string
}

func NewValidateConnector(
	rootPath string,
	schemaPath string,
	connectorPaths []string,
) *ValidateConnector {
	cp := make([]string, len(connectorPaths))
	for i, v := range connectorPaths {
		cp[i] = fmt.Sprintf("%s/%s", rootPath, v)
	}
	return &ValidateConnector{
		rootPath:       rootPath,
		schemaPath:     fmt.Sprintf("%s/%s", rootPath, schemaPath),
		connectorPaths: cp,
	}
}

func (vc *ValidateConnector) Run() {
	var wg sync.WaitGroup
	wg.Add(len(vc.connectorPaths))
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", vc.schemaPath))
	for _, v := range vc.connectorPaths {
		go func(path string) {
			defer wg.Done()
			vc.validateConnector(path, schemaLoader)
		}(v)
	}
	wg.Wait()
	fmt.Println("Validation success")
	os.Exit(Success)
}

func (vc *ValidateConnector) validateConnector(path string, schemaLoader gojsonschema.JSONLoader) {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", path))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("Schema Validation Error", err)
	}

	if !result.Valid() {
		fmt.Printf("The connecot %s is not valid. see errors :\n", path)
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		os.Exit(Failed)
	}
}
