package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cp "github.com/otiai10/copy"
)

func getUserInput() {
	// Get Namespace / App Name
	fmt.Println("Enter the namespace for the app (excluding prod/stage):")
	inputReader := bufio.NewReader(os.Stdin)
	appName, _ := inputReader.ReadString('\n')
	appName = strings.TrimSuffix(appName, "\n")

	// Get Image name
	fmt.Println("Enter docker image name (without tags):")
	inputReader = bufio.NewReader(os.Stdin)
	imageName, _ := inputReader.ReadString('\n')
	imageName = strings.TrimSuffix(imageName, "\n")

	// Get Image name
	fmt.Println("Enter docker image production tag (Eg: production-311e2b7a):")
	inputReader = bufio.NewReader(os.Stdin)
	imageTagProd, _ := inputReader.ReadString('\n')
	imageTagProd = strings.TrimSuffix(imageTagProd, "\n")

	// Get Image name
	fmt.Println("Enter docker image staging tag (Eg: staging-3a8217e2):")
	inputReader = bufio.NewReader(os.Stdin)
	imageTagStage, _ := inputReader.ReadString('\n')
	imageTagStage = strings.TrimSuffix(imageTagStage, "\n")

	// Get Port
	fmt.Println("Enter docker container port (Eg: 3000):")
	inputReader = bufio.NewReader(os.Stdin)
	containerPort, _ := inputReader.ReadString('\n')
	containerPort = strings.TrimSuffix(containerPort, "\n")

	// Get Domain Name: Production
	fmt.Println("Enter Domain name for Production Ingress (Eg: prod.example.com):")
	inputReader = bufio.NewReader(os.Stdin)
	ingressHostProd, _ := inputReader.ReadString('\n')
	ingressHostProd = strings.TrimSuffix(ingressHostProd, "\n")

	// Get Secret Name: Production
	fmt.Println("Enter Secret Key for Production Ingress (Eg: k8s-tls-secret-replica):")
	inputReader = bufio.NewReader(os.Stdin)
	ingressSecretKeyProd, _ := inputReader.ReadString('\n')
	ingressSecretKeyProd = strings.TrimSuffix(ingressSecretKeyProd, "\n")

	// Get Domain Name: Staging
	fmt.Println("Enter Domain name for Staging Ingress (Eg: stage.example.com):")
	inputReader = bufio.NewReader(os.Stdin)
	ingressHostStage, _ := inputReader.ReadString('\n')
	ingressHostStage = strings.TrimSuffix(ingressHostStage, "\n")

	// Get Secret Name: Staging
	fmt.Println("Enter Secret Key for Staging Ingress (Eg: k8s-tls-secret-replica):")
	inputReader = bufio.NewReader(os.Stdin)
	ingressSecretKeyStage, _ := inputReader.ReadString('\n')
	ingressSecretKeyStage = strings.TrimSuffix(ingressSecretKeyStage, "\n")

	replaceText(appName)
	fmt.Println("Your docker image name is:", imageName)
	fmt.Println("Your docker image production tag is:", imageTagProd)
	fmt.Println("Your docker image staging tag is:", imageTagStage)
	fmt.Println("Your Sontainer port is:", containerPort)
	fmt.Println("Your Production domain is:", ingressHostProd)
	fmt.Println("Your Production secret key is:", ingressSecretKeyProd)
	fmt.Println("Your Staging domain is:", ingressHostStage)
	fmt.Println("Your Staging secret key is:", ingressSecretKeyStage)

}

func replaceText(appName string) {
	err := cp.Copy("github.com/pravinbanjade/k8s-config-generator/template", "./generated-k8s-config")
	fmt.Println(err)
	fmt.Println("Your app name is:", appName)
	filePath := "./base/service-account.yaml"
	fileData, _ := ioutil.ReadFile(filePath)
	fileString := string(fileData)
	fileString = strings.ReplaceAll(fileString, "<appName>", appName)
	fileData = []byte(fileString)
	_ = ioutil.WriteFile(filePath, fileData, 0o600)
}

func main() {
	getUserInput()
}
