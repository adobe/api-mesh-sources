package main

import (
	"errors"
	"fmt"
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

func (vc *ValidateConnector) Run() error {
	var wg sync.WaitGroup
	var errList []error
	wg.Add(len(vc.connectorPaths))
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", vc.schemaPath))
	for _, v := range vc.connectorPaths {
		go func(path string) {
			defer wg.Done()
			err := vc.validateConnector(path, schemaLoader)
			if err != nil {
				errList = append(errList, err)
			}
		}(v)
	}
	wg.Wait()
	if len(errList) != 0 {
		err := errors.New("The connectors validation failed:")
		for _, connError := range errList {
			err = fmt.Errorf("%w\n %s\n", err, connError)
		}
		fmt.Printf("%s", err)
		return err
	}

	fmt.Println("Validation success")
	return nil
}

func (vc *ValidateConnector) validateConnector(path string, schemaLoader gojsonschema.JSONLoader) error {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", path))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("Schema Validation Error", err)
	}

	if !result.Valid() {
		err := errors.New(fmt.Sprintf("Connector: %s\n", path))

		for _, desc := range result.Errors() {
			err = fmt.Errorf("%w %s\n", err, desc)
		}
		return err
	}
	return nil
}
