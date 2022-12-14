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

func makeDir_NodeTemplate() {
	err := os.Mkdir(appName, 0755)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir(appName+"/base", 0755)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir(appName+"/base/common", 0755)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir(appName+"/base/webserver", 0755)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir(appName+"/production", 0755)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir(appName+"/staging", 0755)
	if err != nil {
		fmt.Println(err)
	}
}

func importYamlFile_NodeTemplate(yamlUrl string) {
	// Import template YAML file
	resp, err := http.Get(yamlUrl)
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

	replacePlaceholder_NodeTemplate(string(yamlTemplate), yamlUrl)
}

func replacePlaceholder_NodeTemplate(yamlTemplate, yamlUrl string) {
	// Replace placeholders in template with user input
	yaml := strings.Replace(yamlTemplate, "{{appName}}", appName, -1)
	yaml = strings.Replace(yaml, "{{imageName}}", imageName, -1)
	yaml = strings.Replace(yaml, "{{imageTagProd}}", imageTagProd, -1)
	yaml = strings.Replace(yaml, "{{imageTagStage}}", imageTagStage, -1)
	yaml = strings.Replace(yaml, "{{containerPort}}", containerPort, -1)
	yaml = strings.Replace(yaml, "{{ingressHostProd}}", ingressHostProd, -1)
	yaml = strings.Replace(yaml, "{{ingressSecretKeyProd}}", ingressSecretKeyProd, -1)
	yaml = strings.Replace(yaml, "{{ingressHostStage}}", ingressHostStage, -1)
	yaml = strings.Replace(yaml, "{{ingressSecretKeyStage}}", ingressSecretKeyStage, -1)

	filterYamlUrl := strings.Replace(yamlUrl, "https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs", "", -1)

	// Write generated YAML to file
	err := ioutil.WriteFile(appName+filterYamlUrl, []byte(yaml), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

	makeDir_NodeTemplate()

	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/base/common/service-account.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/base/webserver/deployment.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/base/webserver/service.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/base/kustomization.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/base/resource-quota.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/base/vpa.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/production/config.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/production/ingress.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/production/kustomization.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/production/namespace.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/production/secret.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/staging/config.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/staging/ingress.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/staging/kustomization.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/staging/namespace.yaml")
	importYamlFile_NodeTemplate("https://raw.githubusercontent.com/pravinbanjade/k8s-config-generator/main/src/templates/nodejs/staging/secret.yaml")

}

func laravelConfig() {
	fmt.Println("laravel")
}

func main() {
	chooseTemplate()
}
