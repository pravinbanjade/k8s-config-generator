package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func chooseTemplate() {

	// Print the menu options
	pterm.Print("\n")
	pterm.Info.Println("Please select an template to use: (Default node.js)")
	pterm.Print("\n")

	options := []string{"1. node.js (next, nuxt, express, etc.)", "2. Laravel (php-fpm with nginx)"}

	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
	pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedOption))

	switch selectedOption {
	case "1. node.js (next, nuxt, express, etc.)":
		nodeConfig()
	case "2. Laravel (php-fpm with nginx)":
		laravelConfig()
	default:
		nodeConfig()
	}
}

func nodeConfig() {
	fmt.Println("node")

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

func laravelConfig() {
	fmt.Println("laravel")
}

func main() {
	chooseTemplate()
}
