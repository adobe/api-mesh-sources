/*
Copyright 2022 Adobe. All rights reserved.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or implied. See the License for the specific language
governing permissions and limitations under the License.
*/
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
