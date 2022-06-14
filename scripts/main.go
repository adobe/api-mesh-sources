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

type ActionRouter map[string]func() CommandInterface

func main() {
	commandType := os.Args[1]
	router := make(ActionRouter)
	router["validate-connector"] = func() CommandInterface { return NewValidateConnector(os.Args[2], os.Args[3], os.Args[4:]) }
	router["collect-metadata"] = func() CommandInterface { NewCollectMetadata(os.Args[2], os.Args[3], os.Args[4], os.Args[5:]) }
	router[commandType]().Run()
}

func LogError(err error, info interface{}) {
	if err != nil {
		fmt.Printf("%s: %v", info, err.Error())
		os.Exit(Failed)
	}
}
