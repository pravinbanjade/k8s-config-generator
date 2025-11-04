package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Kubernetes resource structs
type Deployment struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Spec       DeploymentSpec    `yaml:"spec"`
}

type Service struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   Metadata   `yaml:"metadata"`
	Spec       ServiceSpec `yaml:"spec"`
}

type Ingress struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       IngressSpec  `yaml:"spec"`
}

type ServiceAccount struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
}

type ConfigMap struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Data       map[string]string `yaml:"data,omitempty"`
}

type Secret struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Type       string            `yaml:"type,omitempty"`
	StringData map[string]string `yaml:"stringData,omitempty"`
	Data       map[string]string `yaml:"data,omitempty"`
}

type Namespace struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}

type ResourceQuota struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Spec       ResourceQuotaSpec `yaml:"spec"`
}

type VPA struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       VPASpec   `yaml:"spec"`
}

type Metadata struct {
	Name        string            `yaml:"name,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type DeploymentSpec struct {
	Replicas *int32            `yaml:"replicas,omitempty"`
	Selector Selector          `yaml:"selector"`
	Strategy DeploymentStrategy `yaml:"strategy,omitempty"`
	Template PodTemplate        `yaml:"template"`
}

type DeploymentStrategy struct {
	Type          string                 `yaml:"type,omitempty"`
	RollingUpdate map[string]interface{} `yaml:"rollingUpdate,omitempty"`
}

type Selector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type PodTemplate struct {
	Metadata Metadata    `yaml:"metadata"`
	Spec     PodSpec     `yaml:"spec"`
}

type PodSpec struct {
	ServiceAccountName string                 `yaml:"serviceAccountName,omitempty"`
	ImagePullSecrets   []ImagePullSecretRef   `yaml:"imagePullSecrets,omitempty"`
	SecurityContext    map[string]interface{} `yaml:"securityContext,omitempty"`
	Containers         []Container            `yaml:"containers"`
	NodeSelector       map[string]string      `yaml:"nodeSelector,omitempty"`
	Affinity           map[string]interface{} `yaml:"affinity,omitempty"`
	Tolerations        []map[string]interface{} `yaml:"tolerations,omitempty"`
}

type ImagePullSecretRef struct {
	Name string `yaml:"name"`
}

type Container struct {
	Name            string                 `yaml:"name"`
	Image           string                 `yaml:"image"`
	ImagePullPolicy string                 `yaml:"imagePullPolicy,omitempty"`
	Ports           []ContainerPort        `yaml:"ports,omitempty"`
	Env             []map[string]string    `yaml:"env,omitempty"`
	EnvFrom         []map[string]interface{} `yaml:"envFrom,omitempty"`
	Resources       map[string]interface{} `yaml:"resources,omitempty"`
	SecurityContext map[string]interface{} `yaml:"securityContext,omitempty"`
	LivenessProbe   map[string]interface{} `yaml:"livenessProbe,omitempty"`
	ReadinessProbe  map[string]interface{} `yaml:"readinessProbe,omitempty"`
}

type ContainerPort struct {
	Name          string `yaml:"name"`
	ContainerPort int    `yaml:"containerPort"`
	Protocol      string `yaml:"protocol,omitempty"`
}

type ServiceSpec struct {
	Type     string            `yaml:"type,omitempty"`
	Ports    []ServicePort     `yaml:"ports"`
	Selector map[string]string `yaml:"selector"`
}

type ServicePort struct {
	Port       int    `yaml:"port"`
	TargetPort string `yaml:"targetPort"`
	Protocol   string `yaml:"protocol,omitempty"`
	Name       string `yaml:"name,omitempty"`
}

type IngressSpec struct {
	IngressClassName string         `yaml:"ingressClassName,omitempty"`
	TLS              []IngressTLS   `yaml:"tls,omitempty"`
	Rules            []IngressRule  `yaml:"rules"`
}

type IngressTLS struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

type IngressRule struct {
	Host string      `yaml:"host"`
	HTTP IngressHTTP `yaml:"http"`
}

type IngressHTTP struct {
	Paths []IngressPath `yaml:"paths"`
}

type IngressPath struct {
	Path     string              `yaml:"path"`
	PathType string              `yaml:"pathType"`
	Backend  IngressPathBackend  `yaml:"backend"`
}

type IngressPathBackend struct {
	Service IngressService `yaml:"service"`
}

type IngressService struct {
	Name string          `yaml:"name"`
	Port IngressServicePort `yaml:"port"`
}

type IngressServicePort struct {
	Number int `yaml:"number"`
}

type ResourceQuotaSpec struct {
	Hard map[string]string `yaml:"hard"`
}

type VPASpec struct {
	TargetRef   VPATargetRef   `yaml:"targetRef"`
	UpdatePolicy map[string]string `yaml:"updatePolicy"`
}

type VPATargetRef struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Name       string `yaml:"name"`
}

var (
	appName              string
	namespace            string
	imageRepo            string
	imageTag             string
	containerPort        int
	env                  string
	allEnvironments      bool
	replicas             int
	ingressEnabled       bool
	ingressHost          string
	ingressClass         string
	ingressTLSSecret     string
	ingressHostStage     string
	ingressHostProd      string
	ingressTLSSecretStage string
	ingressTLSSecretProd  string
	imageTagStage        string
	imageTagProd         string
	imagePullSecrets     []string
	serviceAccount       string
	createSA             bool
	vpaEnabled           bool
	resourceQuotaEnabled bool
	resourceRequestsCPU  string
	resourceRequestsMemory string
	resourceLimitsCPU    string
	resourceLimitsMemory string
	render               bool
	outputDir            string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "k8s-config-generator",
		Short: "Generate Kubernetes manifests",
		Long:  "A tool to generate Kubernetes manifests with CLI flags",
		RunE:  run,
	}

	// Required flags (optional - will prompt if not provided)
	rootCmd.Flags().StringVar(&appName, "app-name", "", "Application name")
	rootCmd.Flags().StringVar(&imageRepo, "image-repo", "", "Docker image repository")
	rootCmd.Flags().StringVar(&imageTag, "image-tag", "", "Docker image tag")

	// Optional flags
	rootCmd.Flags().StringVar(&namespace, "namespace", "", "Kubernetes namespace")
	rootCmd.Flags().IntVar(&containerPort, "container-port", 3000, "Container port")
	rootCmd.Flags().StringVar(&env, "env", "", "Environment (staging|production)")
	rootCmd.Flags().BoolVar(&allEnvironments, "all-environments", false, "Generate manifests for both staging and production")
	rootCmd.Flags().StringVar(&imageTagStage, "image-tag-stage", "", "Docker image tag for staging (used with --all-environments)")
	rootCmd.Flags().StringVar(&imageTagProd, "image-tag-prod", "", "Docker image tag for production (used with --all-environments)")
	rootCmd.Flags().StringVar(&ingressHostStage, "ingress-host-stage", "", "Ingress host for staging")
	rootCmd.Flags().StringVar(&ingressHostProd, "ingress-host-prod", "", "Ingress host for production")
	rootCmd.Flags().StringVar(&ingressTLSSecretStage, "ingress-tls-secret-stage", "", "Ingress TLS secret for staging")
	rootCmd.Flags().StringVar(&ingressTLSSecretProd, "ingress-tls-secret-prod", "", "Ingress TLS secret for production")
	rootCmd.Flags().IntVar(&replicas, "replicas", 1, "Number of replicas")
	rootCmd.Flags().BoolVar(&ingressEnabled, "ingress-enabled", false, "Enable ingress")
	rootCmd.Flags().StringVar(&ingressHost, "ingress-host", "", "Ingress host")
	rootCmd.Flags().StringVar(&ingressClass, "ingress-class", "nginx", "Ingress class name")
	rootCmd.Flags().StringVar(&ingressTLSSecret, "ingress-tls-secret", "", "Ingress TLS secret name")
	rootCmd.Flags().StringArrayVar(&imagePullSecrets, "image-pull-secret", []string{}, "Image pull secret name (can be repeated)")
	rootCmd.Flags().StringVar(&serviceAccount, "service-account", "", "Service account name")
	rootCmd.Flags().BoolVar(&createSA, "create-service-account", true, "Create service account")
	rootCmd.Flags().BoolVar(&vpaEnabled, "vpa-enabled", false, "Enable VPA")
	rootCmd.Flags().BoolVar(&resourceQuotaEnabled, "resource-quota-enabled", false, "Enable resource quota")
	rootCmd.Flags().StringVar(&resourceRequestsCPU, "resources-requests-cpu", "", "Resource requests CPU")
	rootCmd.Flags().StringVar(&resourceRequestsMemory, "resources-requests-memory", "", "Resource requests memory")
	rootCmd.Flags().StringVar(&resourceLimitsCPU, "resources-limits-cpu", "", "Resource limits CPU")
	rootCmd.Flags().StringVar(&resourceLimitsMemory, "resources-limits-memory", "", "Resource limits memory")

	// Output flags
	rootCmd.Flags().BoolVar(&render, "render", false, "Render manifests to stdout")
	rootCmd.Flags().StringVar(&outputDir, "output-dir", "", "Output directory for rendered manifests")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func promptForInputs() error {
	pterm.Print("\n")
	pterm.Info.Println("Welcome to k8s-config-generator!")
	pterm.Info.Println("Please provide the following information:")
	pterm.Print("\n")

	reader := bufio.NewReader(os.Stdin)

	// Required fields
	if appName == "" {
		pterm.Print("Application name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		appName = strings.TrimSpace(input)
		if appName == "" {
			return fmt.Errorf("application name cannot be empty")
		}
	}

	if imageRepo == "" {
		pterm.Print("Docker image repository (e.g., registry.example.com/myapp): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		imageRepo = strings.TrimSpace(input)
		if imageRepo == "" {
			return fmt.Errorf("image repository cannot be empty")
		}
	}

	// Optional fields
	if namespace == "" {
		pterm.Print("Kubernetes namespace (press Enter to skip): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		namespace = strings.TrimSpace(input)
	}

	// Ask for environment first, then image tags based on selection
	if env == "" {
		options := []string{"staging", "production", "both (staging & production)", "skip"}
		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Environment:")
		if selectedOption == "both (staging & production)" {
			allEnvironments = true
			// Prompt for staging-specific image tag
			if imageTagStage == "" {
				pterm.Print("Docker image tag for staging (e.g., staging-123): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}
				imageTagStage = strings.TrimSpace(input)
				if imageTagStage == "" {
					return fmt.Errorf("staging image tag cannot be empty")
				}
			}
			// Prompt for production-specific image tag
			if imageTagProd == "" {
				pterm.Print("Docker image tag for production (e.g., production-456): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}
				imageTagProd = strings.TrimSpace(input)
				if imageTagProd == "" {
					return fmt.Errorf("production image tag cannot be empty")
				}
			}
		} else if selectedOption != "skip" {
			env = selectedOption
			// Ask for image tag for single environment
			if imageTag == "" {
				pterm.Print("Docker image tag (e.g., v1.0.0 or staging-123): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}
				imageTag = strings.TrimSpace(input)
				if imageTag == "" {
					return fmt.Errorf("image tag cannot be empty")
				}
			}
		}
	} else {
		// Environment was provided via flag, but image tag might not be
		if imageTag == "" && !allEnvironments {
			pterm.Print("Docker image tag (e.g., v1.0.0 or staging-123): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			imageTag = strings.TrimSpace(input)
			if imageTag == "" {
				return fmt.Errorf("image tag cannot be empty")
			}
		}
	}

	if containerPort == 3000 {
		pterm.Print("Container port (default: 3000, press Enter for default): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.TrimSpace(input)
		if input != "" {
			port, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid port number: %w", err)
			}
			containerPort = port
		}
	}

	if replicas == 1 {
		pterm.Print("Number of replicas (default: 1, press Enter for default): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.TrimSpace(input)
		if input != "" {
			rep, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid replica count: %w", err)
			}
			replicas = rep
		}
	}

	// Ingress configuration
	if !ingressEnabled {
		options := []string{"yes", "no"}
		selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Enable ingress? (yes/no):")
		if selectedOption == "yes" {
			ingressEnabled = true

			if allEnvironments {
				// Prompt for both staging and production ingress hosts
				if ingressHostStage == "" {
					pterm.Print("Ingress host for staging (e.g., stage.example.com): ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					ingressHostStage = strings.TrimSpace(input)
				}
				if ingressHostProd == "" {
					pterm.Print("Ingress host for production (e.g., prod.example.com): ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					ingressHostProd = strings.TrimSpace(input)
				}
				if ingressTLSSecretStage == "" {
					pterm.Print("Ingress TLS secret for staging (press Enter to skip): ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					ingressTLSSecretStage = strings.TrimSpace(input)
				}
				if ingressTLSSecretProd == "" {
					pterm.Print("Ingress TLS secret for production (press Enter to skip): ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					ingressTLSSecretProd = strings.TrimSpace(input)
				}
			} else {
				// Single environment - prompt for regular ingress host
				if ingressHost == "" {
					pterm.Print("Ingress host (e.g., app.example.com): ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					ingressHost = strings.TrimSpace(input)
				}
				if ingressTLSSecret == "" {
					pterm.Print("Ingress TLS secret name (press Enter to skip): ")
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					ingressTLSSecret = strings.TrimSpace(input)
				}
			}

			if ingressClass == "nginx" {
				pterm.Print("Ingress class name (default: nginx, press Enter for default): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}
				input = strings.TrimSpace(input)
				if input != "" {
					ingressClass = input
				}
			}
		}
	}

	// Default to create files if no output mode specified
	// We'll create files unless explicitly told otherwise

	pterm.Print("\n")
	pterm.Success.Println("Configuration collected!")
	pterm.Print("\n")

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	// Check if required values are provided
	// If any required value is missing, prompt interactively
	requiredFlagsProvided := appName != "" && imageRepo != "" && imageTag != ""

	// If required flags not provided, prompt for input interactively
	if !requiredFlagsProvided {
		if err := promptForInputs(); err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}
	}

	// Validate required values after prompting
	if appName == "" {
		return fmt.Errorf("application name is required")
	}
	if imageRepo == "" {
		return fmt.Errorf("image repository is required")
	}
	
	// Validate image tags based on environment selection
	if allEnvironments {
		if imageTagStage == "" {
			return fmt.Errorf("staging image tag is required (use --image-tag-stage or provide in interactive mode)")
		}
		if imageTagProd == "" {
			return fmt.Errorf("production image tag is required (use --image-tag-prod or provide in interactive mode)")
		}
	} else {
		if imageTag == "" {
			return fmt.Errorf("image tag is required")
		}
	}

	// Handle different output modes
	if render {
		if outputDir != "" {
			return generateManifestsToDir(outputDir)
		}
		// Default: render to stdout
		return generateManifestsToStdout()
	}

	// Default behavior: create folder structure with files
	if !allEnvironments {
		return createManifestFiles(env)
	}

	return createManifestFilesForAllEnvironments()
}

// Generate manifests directly to stdout
func generateManifestsToStdout() error {
	manifests, err := generateManifests(env, imageTag, namespace, ingressHost, ingressTLSSecret)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		fmt.Println("---")
		data, err := yaml.Marshal(manifest)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
		fmt.Print(string(data))
	}
	return nil
}

// Generate manifests to output directory
func generateManifestsToDir(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	manifests, err := generateManifests(env, imageTag, namespace, ingressHost, ingressTLSSecret)
	if err != nil {
		return err
	}

	var filesCreated []string
	for _, manifest := range manifests {
		filename, err := writeManifestToFile(manifest, outputDir)
		if err != nil {
			return err
		}
		filesCreated = append(filesCreated, filename)
	}

	pterm.Success.Printf("Created %d manifest files in directory: %s\n", len(filesCreated), outputDir)
	return nil
}

// Create manifest files in app directory
func createManifestFiles(envName string) error {
	// Determine values based on environment
	tag := imageTag
	ns := namespace
	host := ingressHost
	tlsSecret := ingressTLSSecret

	// Use environment-specific namespace if not provided
	if ns == "" && envName != "" {
		ns = fmt.Sprintf("%s-%s", appName, envName)
	}

	manifests, err := generateManifests(envName, tag, ns, host, tlsSecret)
	if err != nil {
		return err
	}

	// Create output directory
	outputDir := appName
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var filesCreated []string
	for _, manifest := range manifests {
		filename, err := writeManifestToFile(manifest, outputDir)
		if err != nil {
			return err
		}
		filesCreated = append(filesCreated, filename)
	}

	pterm.Success.Printf("Created %d manifest files in directory: %s\n", len(filesCreated), outputDir)
	pterm.Info.Println("Files created:")
	for _, file := range filesCreated {
		pterm.Printf("  - %s\n", filepath.Join(outputDir, file))
	}

	return nil
}

// Create manifest files for all environments
func createManifestFilesForAllEnvironments() error {
	environments := []struct {
		name         string
		imageTag     string
		namespace    string
		ingressHost  string
		tlsSecret    string
	}{
		{
			name:         "staging",
			imageTag:     imageTagStage,
			namespace:    fmt.Sprintf("%s-staging", appName),
			ingressHost:  ingressHostStage,
			tlsSecret:    ingressTLSSecretStage,
		},
		{
			name:         "production",
			imageTag:     imageTagProd,
			namespace:    fmt.Sprintf("%s-production", appName),
			ingressHost:  ingressHostProd,
			tlsSecret:    ingressTLSSecretProd,
		},
	}

	// Create base directory
	baseDir := appName
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	pterm.Info.Printf("Generating manifests for all environments in: %s\n", baseDir)
	pterm.Print("\n")

	// Generate manifests for each environment
	for _, envConfig := range environments {
		pterm.Info.Printf("Generating %s environment...\n", envConfig.name)

		manifests, err := generateManifests(envConfig.name, envConfig.imageTag, envConfig.namespace, envConfig.ingressHost, envConfig.tlsSecret)
		if err != nil {
			return fmt.Errorf("failed to generate manifests for %s: %w", envConfig.name, err)
		}

		// Create environment-specific directory
		envDir := filepath.Join(baseDir, envConfig.name)
		if err := os.MkdirAll(envDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", envConfig.name, err)
		}

		// Write manifests to environment directory
		var filesCreated []string
		for _, manifest := range manifests {
			filename, err := writeManifestToFile(manifest, envDir)
			if err != nil {
				return fmt.Errorf("failed to write manifest for %s: %w", envConfig.name, err)
			}
			filesCreated = append(filesCreated, filename)
		}

		pterm.Success.Printf("  ✓ %s environment manifests created in %s (%d files)\n", envConfig.name, envDir, len(filesCreated))
	}

	pterm.Print("\n")
	pterm.Success.Printf("All environments generated successfully!\n")
	pterm.Info.Printf("Directory structure:\n")
	pterm.Printf("  %s/\n", baseDir)
	pterm.Printf("    ├── staging/\n")
	pterm.Printf("    └── production/\n")

	return nil
}

// Generate all Kubernetes manifests
func generateManifests(envName, tag, ns, ingressHostVal, tlsSecretVal string) ([]interface{}, error) {
	var manifests []interface{}

	// Determine if we should enable ConfigMap and Secret
	enableConfigMap := envName == "staging" || envName == "production"
	enableIngress := ingressEnabled && (ingressHostVal != "" || allEnvironments)

	// Namespace
	if ns != "" {
		namespace := createNamespace(ns)
		manifests = append(manifests, namespace)
	}

	// ServiceAccount
	if createSA {
		saName := appName
		if serviceAccount != "" {
			saName = serviceAccount
		}
		sa := createServiceAccount(saName, ns)
		manifests = append(manifests, sa)
	}

	// ConfigMap
	if enableConfigMap {
		configMap := createConfigMap(envName, ns)
		manifests = append(manifests, configMap)
	}

	// Secret
	if enableConfigMap {
		secret := createSecret(envName, ns)
		manifests = append(manifests, secret)
	}

	// Deployment
	deployment := createDeployment(tag, ns, enableConfigMap)
	manifests = append(manifests, deployment)

	// Service
	service := createService(ns)
	manifests = append(manifests, service)

	// Ingress
	if enableIngress && ingressHostVal != "" {
		ingress := createIngress(ingressHostVal, tlsSecretVal, ns)
		manifests = append(manifests, ingress)
	}

	// ResourceQuota
	if resourceQuotaEnabled {
		quota := createResourceQuota(ns)
		manifests = append(manifests, quota)
	}

	// VPA
	if vpaEnabled {
		vpa := createVPA(ns)
		manifests = append(manifests, vpa)
	}

	return manifests, nil
}

// Create Namespace resource
func createNamespace(ns string) *Namespace {
	return &Namespace{
		APIVersion: "v1",
		Kind:       "Namespace",
		Metadata: Metadata{
			Name: ns,
		},
	}
}

// Create ServiceAccount resource
func createServiceAccount(name, ns string) *ServiceAccount {
	return &ServiceAccount{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
		Metadata: Metadata{
			Name:      name,
			Namespace: ns,
		},
	}
}

// Create ConfigMap resource
func createConfigMap(envName, ns string) *ConfigMap {
	data := map[string]string{
		"APP_ENV": envName,
	}
	return &ConfigMap{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: Metadata{
			Name:      appName,
			Namespace: ns,
		},
		Data: data,
	}
}

// Create Secret resource
func createSecret(envName, ns string) *Secret {
	stringData := map[string]string{
		"APP_KEY": "",
	}
	return &Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Metadata: Metadata{
			Name:      appName,
			Namespace: ns,
		},
		Type:       "Opaque",
		StringData: stringData,
	}
}

// Create Deployment resource
func createDeployment(tag, ns string, useEnvFrom bool) *Deployment {
	// Selector labels are required for Deployment selector and pod template
	selectorLabels := map[string]string{
		"app.kubernetes.io/name":     "k8s-config-generator",
		"app.kubernetes.io/instance": appName,
		"tier":                       "webserver",
		"layer":                      "node",
	}

	replicasInt32 := int32(replicas)

	// Build image pull secrets
	var imagePullSecretsRefs []ImagePullSecretRef
	for _, secret := range imagePullSecrets {
		imagePullSecretsRefs = append(imagePullSecretsRefs, ImagePullSecretRef{Name: secret})
	}

	// Build container
	container := Container{
		Name:            fmt.Sprintf("%s-node", appName),
		Image:           fmt.Sprintf("%s:%s", imageRepo, tag),
		ImagePullPolicy: "IfNotPresent",
		Ports: []ContainerPort{
			{
				Name:          "http-port",
				ContainerPort: containerPort,
				Protocol:      "TCP",
			},
		},
		SecurityContext: map[string]interface{}{
			"allowPrivilegeEscalation": false,
			"capabilities": map[string]interface{}{
				"drop": []string{"ALL"},
			},
			"privileged": false,
		},
	}

	// Add envFrom if ConfigMap/Secret are enabled
	if useEnvFrom {
		container.EnvFrom = []map[string]interface{}{
			{
				"configMapRef": map[string]string{
					"name": appName,
				},
			},
			{
				"secretRef": map[string]string{
					"name": appName,
				},
			},
		}
	}

	// Add resources if provided
	if resourceRequestsCPU != "" || resourceRequestsMemory != "" || resourceLimitsCPU != "" || resourceLimitsMemory != "" {
		resources := make(map[string]interface{})
		requests := make(map[string]string)
		limits := make(map[string]string)

		if resourceRequestsCPU != "" {
			requests["cpu"] = resourceRequestsCPU
		}
		if resourceRequestsMemory != "" {
			requests["memory"] = resourceRequestsMemory
		}
		if resourceLimitsCPU != "" {
			limits["cpu"] = resourceLimitsCPU
		}
		if resourceLimitsMemory != "" {
			limits["memory"] = resourceLimitsMemory
		}

		if len(requests) > 0 {
			resources["requests"] = requests
		}
		if len(limits) > 0 {
			resources["limits"] = limits
		}
		container.Resources = resources
	}

	saName := appName
	if serviceAccount != "" {
		saName = serviceAccount
	}

	deployment := &Deployment{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: Metadata{
			Name:      fmt.Sprintf("%s-node", appName),
			Namespace: ns,
		},
		Spec: DeploymentSpec{
			Replicas: &replicasInt32,
			Selector: Selector{
				MatchLabels: selectorLabels,
			},
			Strategy: DeploymentStrategy{
				Type: "RollingUpdate",
				RollingUpdate: map[string]interface{}{
					"maxSurge":       "50%",
					"maxUnavailable": "25%",
				},
			},
			Template: PodTemplate{
				Metadata: Metadata{
					Labels: selectorLabels,
				},
				Spec: PodSpec{
					ServiceAccountName: saName,
					ImagePullSecrets:   imagePullSecretsRefs,
					Containers:         []Container{container},
				},
			},
		},
	}

	return deployment
}

// Create Service resource
func createService(ns string) *Service {
	// Selector labels are required for Service selector
	selectorLabels := map[string]string{
		"app.kubernetes.io/name":     "k8s-config-generator",
		"app.kubernetes.io/instance": appName,
		"tier":                       "webserver",
		"layer":                      "node",
	}

	return &Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:      appName,
			Namespace: ns,
		},
		Spec: ServiceSpec{
			Type: "ClusterIP",
			Ports: []ServicePort{
				{
					Port:       80,
					TargetPort: "http-port",
					Protocol:   "TCP",
					Name:       "http",
				},
			},
			Selector: selectorLabels,
		},
	}
}

// Create Ingress resource
func createIngress(host, tlsSecret, ns string) *Ingress {
	ingress := &Ingress{
		APIVersion: "networking.k8s.io/v1",
		Kind:       "Ingress",
		Metadata: Metadata{
			Name:      fmt.Sprintf("%s-ingress", appName),
			Namespace: ns,
		},
		Spec: IngressSpec{
			IngressClassName: ingressClass,
			Rules: []IngressRule{
				{
					Host: host,
					HTTP: IngressHTTP{
						Paths: []IngressPath{
							{
								Path:     "/",
								PathType: "Prefix",
								Backend: IngressPathBackend{
									Service: IngressService{
										Name: appName,
										Port: IngressServicePort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Add TLS if secret provided
	if tlsSecret != "" {
		ingress.Spec.TLS = []IngressTLS{
			{
				Hosts:      []string{host},
				SecretName: tlsSecret,
			},
		}
	}

	return ingress
}

// Create ResourceQuota resource
func createResourceQuota(ns string) *ResourceQuota {
	hard := map[string]string{
		"limits.cpu":    "200m",
		"limits.memory": "512Mi",
		"requests.cpu":  "200m",
		"requests.memory": "512Mi",
	}

	return &ResourceQuota{
		APIVersion: "v1",
		Kind:       "ResourceQuota",
		Metadata: Metadata{
			Name:      fmt.Sprintf("%s-quota", appName),
			Namespace: ns,
		},
		Spec: ResourceQuotaSpec{
			Hard: hard,
		},
	}
}

// Create VPA resource
func createVPA(ns string) *VPA {
	return &VPA{
		APIVersion: "autoscaling.k8s.io/v1beta2",
		Kind:       "VerticalPodAutoscaler",
		Metadata: Metadata{
			Name:      fmt.Sprintf("%s-node-vpa", appName),
			Namespace: ns,
		},
		Spec: VPASpec{
			TargetRef: VPATargetRef{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       fmt.Sprintf("%s-node", appName),
			},
			UpdatePolicy: map[string]string{
				"updateMode": "Auto",
			},
		},
	}
}

// Write manifest to file and return filename
func writeManifestToFile(manifest interface{}, outputDir string) (string, error) {
	// Get kind and name from manifest
	var kind, name string
	
	switch m := manifest.(type) {
	case *Namespace:
		kind = "namespace"
		name = m.Metadata.Name
	case *ServiceAccount:
		kind = "serviceaccount"
		name = m.Metadata.Name
	case *ConfigMap:
		kind = "configmap"
		name = m.Metadata.Name
	case *Secret:
		kind = "secret"
		name = m.Metadata.Name
	case *Deployment:
		kind = "deployment"
		name = m.Metadata.Name
	case *Service:
		kind = "service"
		name = m.Metadata.Name
	case *Ingress:
		kind = "ingress"
		name = m.Metadata.Name
	case *ResourceQuota:
		kind = "resourcequota"
		name = m.Metadata.Name
	case *VPA:
		kind = "vpa"
		name = m.Metadata.Name
	default:
		kind = "manifest"
		name = "unknown"
	}

	// Generate filename
	cleanName := strings.ToLower(name)
	cleanName = strings.ReplaceAll(cleanName, "/", "-")
	cleanName = strings.ReplaceAll(cleanName, ":", "-")
	filename := fmt.Sprintf("%s-%s.yaml", kind, cleanName)

	filePath := filepath.Join(outputDir, filename)

	// Marshal to YAML
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return filename, nil
}
