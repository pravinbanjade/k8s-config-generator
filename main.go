package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// Import template YAML file
	yamlTemplate, err := ioutil.ReadFile("https://github.com/pravinbanjade/k8s-config-generator/template.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get user input
	fmt.Println("Enter service name:")
	var serviceName string
	fmt.Scanln(&serviceName)

	fmt.Println("Enter app name:")
	var appName string
	fmt.Scanln(&appName)

	fmt.Println("Enter target port:")
	var targetPort string
	fmt.Scanln(&targetPort)

	// Replace placeholders in template with user input
	yaml := strings.Replace(string(yamlTemplate), "{{serviceName}}", serviceName, -1)
	yaml = strings.Replace(yaml, "{{appName}}", appName, -1)
	yaml = strings.Replace(yaml, "{{targetPort}}", targetPort, -1)

	// Write generated YAML to file
	err = ioutil.WriteFile("service.yaml", []byte(yaml), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
