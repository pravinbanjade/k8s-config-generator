package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getUserInput() {

	var appName, imageName, imageTagProd, imageTagStage, containerPort, ingressHostProd, ingressSecretKeyProd, ingressHostStage, ingressSecretKeyStage string
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

func copyDir(src, dest string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			// Create the destination directory
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				return err
			}

			// Recursively copy the contents of the directory
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy the file to the destination
			data, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(destPath, data, file.Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}

func replaceText(appName string) {
	// Copy the "src" directory to the "dest" directory
	err := copyDir("src/template", "src/generated-template")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	getUserInput()
}
