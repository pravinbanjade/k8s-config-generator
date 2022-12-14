package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

var appName, imageName, imageTagProd, imageTagStage, containerPort, ingressHostProd, ingressSecretKeyProd, ingressHostStage, ingressSecretKeyStage string

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

func getUserInput() {

	// Get Namespace / App Name
	fmt.Println("Enter the namespace for the app (excluding prod/stage):")
	fmt.Scan(&appName)

	// Get Image name
	fmt.Println("Enter docker image name (without tags):")
	fmt.Scan(&imageName)

	// Get Image name
	fmt.Println("Enter docker image production tag (Eg: production-311e2b7a):")
	fmt.Scan(&imageTagProd)

	// Get Image name
	fmt.Println("Enter docker image staging tag (Eg: staging-3a8217e2):")
	fmt.Scan(&imageTagStage)

	// Get Port
	fmt.Println("Enter docker container port (Eg: 3000):")
	fmt.Scan(&containerPort)

	// Get Domain Name: Production
	fmt.Println("Enter Domain name for Production Ingress (Eg: prod.example.com):")
	fmt.Scan(&ingressHostProd)

	// Get Secret Name: Production
	fmt.Println("Enter Secret Key for Production Ingress (Eg: k8s-tls-secret-replica):")
	fmt.Scan(&ingressSecretKeyProd)

	// Get Domain Name: Staging
	fmt.Println("Enter Domain name for Staging Ingress (Eg: stage.example.com):")
	fmt.Scan(&ingressHostStage)

	// Get Secret Name: Staging
	fmt.Println("Enter Secret Key for Staging Ingress (Eg: k8s-tls-secret-replica):")
	fmt.Scan(&ingressSecretKeyStage)

}

func nodeConfig() {
	fmt.Println("node")

	// Get user input
	getUserInput()

	fmt.Println("Your app name is:", appName)
	fmt.Println("Your docker image name is:", imageName)
	fmt.Println("Your docker image production tag is:", imageTagProd)
	fmt.Println("Your docker image staging tag is:", imageTagStage)
	fmt.Println("Your Container port is:", containerPort)
	fmt.Println("Your Production domain is:", ingressHostProd)
	fmt.Println("Your Production secret key is:", ingressSecretKeyProd)
	fmt.Println("Your Staging domain is:", ingressHostStage)
	fmt.Println("Your Staging secret key is:", ingressSecretKeyStage)

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
	yaml := strings.Replace(string(yamlTemplate), "{{serviceName}}", appName, -1)
	yaml = strings.Replace(yaml, "{{appName}}", appName, -1)
	yaml = strings.Replace(yaml, "{{targetPort}}", containerPort, -1)

	err = os.Mkdir(appName, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// Write generated YAML to file
	err = ioutil.WriteFile(appName+"/service.yaml", []byte(yaml), 0644)
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
