package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-playground/validator/v10"
)

type Connector struct {
	Name        string      `json:"name" validate:"required"`
	Version     string      `json:"version" validate:"required,semver"`
	Description string      `json:"description" validate:"required"`
	Author      string      `json:"author" validate:"required"`
	Provider    interface{} `json:"provider" validate:"required"`
}

func main() {
	args := os.Args[1:]
	validate := validator.New()
	fmt.Println(args)

	if len(args) != 0 {
		for _, v := range args {
			file, _ := ioutil.ReadFile(v)
			connector := &Connector{}
			err := json.Unmarshal(file, connector)
			if err != nil {
				fmt.Errorf(err.Error())
			}
			err = validate.Struct(connector)
			if err != nil {
				validationErrors := err.(validator.ValidationErrors)
				fmt.Println("err", validationErrors)
			}
		}
	}
	fmt.Println("Validation success")
}
