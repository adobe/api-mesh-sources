package main

import (
	"fmt"
	"os"
)

const (
	Success int = iota
	Failed
)

type CommandInterface interface {
	Run()
}

type ActionRouter map[string]CommandInterface

func main() {
	commandType := os.Args[1]
	router := make(ActionRouter)
	router["validate-connector"] = NewValidateConnector(os.Args[2], os.Args[3], os.Args[4:])
	router["collect-metadata"] = NewCollectMetadata(os.Args[2], os.Args[3], os.Args[4], os.Args[5:])
	router[commandType].Run()
}

func LogError(err error, info interface{}) {
	if err != nil {
		fmt.Printf("%s: %v", info, err.Error())
		os.Exit(Failed)
	}
}
