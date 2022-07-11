package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"golang.org/x/mod/semver"
)

var validate *validator.Validate

type Connector struct {
	Name        string      `json:"name" validate:"required"`
	Version     string      `json:"version" validate:"required,semver"`
	Description string      `json:"description" validate:"required"`
	Author      string      `json:"author" validate:"required"`
	Provider    interface{} `json:"provider" validate:"required"`
}

type ConnectorMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Latest      string   `json:"latest"`
	Versions    []string `json:"versions"`
	Url         string   `json:"url"`
}

type ConnectorsMetadata map[string]*ConnectorMetadata
type CollectMetadata struct {
	rootPath       string
	metaPath       string
	storagePath    string
	connectorPaths []string
}

func NewCollectMetadata(
	rootPath string,
	metadataFilePath string,
	storageFolderPath string,
	connectorPaths []string,
) *CollectMetadata {
	connectorPathsFormatted := make([]string, len(connectorPaths))
	for i, v := range connectorPaths {
		connectorPathsFormatted[i] = fmt.Sprintf("%s/%s", rootPath, v)
	}
	fmt.Println("Test")
	return &CollectMetadata{
		rootPath:       rootPath,
		metaPath:       fmt.Sprintf("%s/%s", rootPath, metadataFilePath),
		storagePath:    fmt.Sprintf("%s/%s", rootPath, storageFolderPath),
		connectorPaths: connectorPathsFormatted,
	}
}

func (cm *CollectMetadata) Run() error {
	connectorsMetadata := make(ConnectorsMetadata)
	connectorsMap, err := ioutil.ReadFile(cm.metaPath)

	if err != nil {
		return fmt.Errorf("%s: %v", "Error on data file", err.Error())
	}

	err = json.Unmarshal(connectorsMap, &connectorsMetadata)
	if err != nil {
		return fmt.Errorf("%s: %v", "Error on unmarshaling connectors metadata", err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(len(cm.connectorPaths))
	for _, v := range cm.connectorPaths {
		go func(cp string, csm ConnectorsMetadata) {
			defer wg.Done()
			cm.processConnector(cp, csm)
		}(v, connectorsMetadata)
	}
	wg.Wait()
	connectorsMapString, err := json.Marshal(connectorsMetadata)
	if err != nil {
		return fmt.Errorf("%s: %v", "Error on marshaling file", err.Error())
	}
	fmt.Println(string(connectorsMapString))
	ioutil.WriteFile(cm.metaPath, connectorsMapString, 0644)
	return nil
}

func (cm *CollectMetadata) processConnector(cp string, csm ConnectorsMetadata) error {
	connector, file, err := cm.getConnector(cp)
	if err != nil {
		return err
	}
	ckn := strings.ToLower(strings.Replace(connector.Name, " ", "-", -1))
	csm[ckn] = cm.getUpdatedMetadata(csm[ckn], connector)
	ioutil.WriteFile(fmt.Sprintf("%s/%s-%s.json", cm.storagePath, connector.Version, ckn), file, 0644)
	return nil
}

func (cm *CollectMetadata) getConnector(path string) (*Connector, []byte, error) {
	validate := validator.New()
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, nil, fmt.Errorf("%s: %v", "Error on reading file", err.Error())
	}
	connector := &Connector{}
	err = json.Unmarshal(file, connector)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %v", "Error on unmarshaling JSON", err.Error())
	}
	err = validate.Struct(connector)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %v", validationErrors, err.Error())
		}
	}
	return connector, file, nil
}

func (cm *CollectMetadata) getUpdatedMetadata(prevMeta *ConnectorMetadata, connector *Connector) *ConnectorMetadata {
	var c *ConnectorMetadata
	if prevMeta != nil {
		c = prevMeta
		compare := semver.Compare(fmt.Sprintf("v%s", c.Latest), fmt.Sprintf("v%s", connector.Version))
		if compare == -1 {
			ckn := strings.ToLower(strings.Replace(connector.Name, " ", "-", -1))
			c.Url = fmt.Sprintf("%s/%s-%s.json", cm.storagePath, connector.Version, ckn)
			c.Versions = append(c.Versions, connector.Version)
			c.Latest = connector.Version
			c.Author = connector.Author
			c.Description = connector.Description
		}
		if compare == 0 {
			c.Author = connector.Author
			c.Description = connector.Description
		}
	} else {
		c = &ConnectorMetadata{
			Name:        connector.Name,
			Latest:      connector.Version,
			Versions:    []string{connector.Version},
			Author:      connector.Author,
			Description: connector.Description,
		}
	}

	return c
}
