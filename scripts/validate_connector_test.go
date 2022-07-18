package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidateConnectorWithInvalidConfiguration(t *testing.T) {
	vc := NewValidateConnector("../", "connector.schema.json", []string{"scripts/test/mock/connectors/invalid_connector.json"})
	err := vc.Run()
	if err == nil {
		t.Errorf("Failed: The mock connector doesn't have required fields: Description")
	}

	errorMock := errors.New("The connectors validation failed:")
	errorMock = fmt.Errorf("%w\n Connector: ..//scripts/test/mock/connectors/invalid_connector.json", errorMock)
	errorMock = fmt.Errorf("%w\n (root): description is required\n\n", errorMock)

	if err.Error() != errorMock.Error() {
		t.Errorf("The error \n %s doesn't equal mock error \n%s", err.Error(), errorMock.Error())
	}
}

func TestValidateConnectorWithValidConfiguration(t *testing.T) {
	vc := NewValidateConnector("../", "connector.schema.json", []string{"scripts/test/mock/connectors/valid_connector.json"})
	err := vc.Run()
	if err != nil {
		t.Errorf("Failed: The mock connector is valid, but validation failed")
	}
}
