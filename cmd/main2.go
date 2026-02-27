package main

import (
	"os"

	apiv1 "github.com/kyma-project/rt-bootstrapper/pkg/api/v1"
)

// write a main function that prints "Hello, World!" to the console

func readConfig2(name string) (*apiv1.Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return apiv1.NewConfig(file)
}

func main() {
	cfg, err := readConfig2("/Users/i316752/SAPDevelop/inner_source/frog_documentation/howto/scripts/install-kim-chart/rt-bootstrapper-config/config2.yaml")
	if err != nil {
		os.Exit(1)
	}

	println("Hello, World!", cfg)
}
