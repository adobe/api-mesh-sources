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
	cp := make([]string, len(connectorPaths))
	for i, v := range connectorPaths {
		cp[i] = fmt.Sprintf("%s/%s", rootPath, v)
	}

	return &CollectMetadata{
		rootPath:       rootPath,
		metaPath:       fmt.Sprintf("%s/%s", rootPath, metadataFilePath),
		storagePath:    fmt.Sprintf("%s/%s", rootPath, storageFolderPath),
		connectorPaths: cp,
	}
}

func (cm *CollectMetadata) Run() {
	connectorsMetadata := make(ConnectorsMetadata)
	connectorsMap, err := ioutil.ReadFile(cm.metaPath)

	LogError(err, "Error on data file")
	err = json.Unmarshal(connectorsMap, &connectorsMetadata)
	LogError(err, "Error on unmarshaling connectors metadata")

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
	LogError(err, "Error on marshaling file")
	ioutil.WriteFile(cm.metaPath, connectorsMapString, 0644)
}

func (cm *CollectMetadata) processConnector(cp string, csm ConnectorsMetadata) {
	connector, file := cm.getConnector(cp)
	ckn := strings.ToLower(strings.Replace(connector.Name, " ", "-", -1))
	csm[ckn] = cm.getUpdatedMetadata(csm[ckn], connector)
	ioutil.WriteFile(fmt.Sprintf("%s/%s-%s.json", cm.storagePath, connector.Version, ckn), file, 0644)
}

func (cm *CollectMetadata) getConnector(path string) (*Connector, []byte) {
	validate := validator.New()
	file, err := ioutil.ReadFile(path)
	LogError(err, "Error on reading file")
	connector := &Connector{}
	err = json.Unmarshal(file, connector)
	LogError(err, "Error on unmarshaling JSON")
	err = validate.Struct(connector)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		LogError(err, validationErrors)
	}
	return connector, file
}

func (cm *CollectMetadata) getUpdatedMetadata(prevMeta *ConnectorMetadata, connector *Connector) *ConnectorMetadata {
	var c *ConnectorMetadata
	if prevMeta != nil {
		c = prevMeta
		if semver.Compare(fmt.Sprintf("v%s", c.Latest), fmt.Sprintf("v%s", connector.Version)) == -1 {
			ckn := strings.ToLower(strings.Replace(connector.Name, " ", "-", -1))
			c.Url = fmt.Sprintf("%s/%s-%s.json", cm.storagePath, connector.Version, ckn)
			c.Versions = append(c.Versions, connector.Version)
			c.Latest = connector.Version
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
