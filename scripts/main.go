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
	Run() error
}

type ActionRouter map[string]func() CommandInterface

func main() {
	commandType := os.Args[1]
	router := make(ActionRouter)
	router["validate-connector"] = func() CommandInterface { return NewValidateConnector(os.Args[2], os.Args[3], os.Args[4:]) }
	router["collect-metadata"] = func() CommandInterface { return NewCollectMetadata(os.Args[2], os.Args[3], os.Args[4], os.Args[5:]) }
	err := router[commandType]().Run()
	if err != nil {
		fmt.Printf("%v", err.Error())
		os.Exit(Failed)
	}
	os.Exit(Success)
}
