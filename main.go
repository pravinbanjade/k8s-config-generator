package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Import template YAML file
	resp, err := http.Get("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/template.yaml")
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	yamlTemplate, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
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
