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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestCollectMetadataAndArchiveCreation(t *testing.T) {
	collectMetadata := NewCollectMetadata(
		"../",
		"scripts/test/mock/connectors/connectors-metadata.json",
		"archive",
		[]string{"scripts/test/mock/connectors/valid_connector.json"},
	)
	err := collectMetadata.Run()
	if err != nil {
		t.Errorf("Failed: Execution error %s", err.Error())
	}

	if _, err := os.Stat("../archive/0.0.1-valid-connector.json"); err != nil && errors.Is(err, os.ErrNotExist) {
		t.Errorf("Failed: Archive file for mock connector doesn't exist, %s", err.Error())
	}

	metadataFile, err := ioutil.ReadFile("../scripts/test/mock/connectors/connectors-metadata.json")
	if err != nil {
		t.Errorf("Failed: Cannot read metadata file, %s", err.Error())
	}

	connectorsMetadata := make(ConnectorsMetadata)
	err = json.Unmarshal(metadataFile, &connectorsMetadata)
	if err != nil {
		t.Errorf("%s: %s", "Error on unmarshaling mock metadata", err.Error())
	}

	if _, ok := connectorsMetadata["valid-connector"]; !ok {
		t.Errorf("Metadata for mock connector doesn't exists in connectors-metadata mock")
	}

	rollbackTestCollectMetadataAndArchiveCreation(t)

}

func rollbackTestCollectMetadataAndArchiveCreation(t *testing.T) {
	err := os.Remove("../archive/0.0.1-valid-connector.json")
	if err != nil {
		t.Errorf("Failed: cannot rollback archive file, %s", err.Error())
	}
	err = ioutil.WriteFile("../scripts/test/mock/connectors/connectors-metadata.json", []byte("{}"), 0644)
	if err != nil {
		t.Errorf("Failed: cannot rollback metadata file, %s", err.Error())
	}
}
