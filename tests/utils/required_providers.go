package utils

import (
	"io"
	"os"
	"text/template"
)

// RequiredProvidersData is the data struct for the Terraform required providers block.
// It should ordinarily be generated using utils.NewRequiredProvidersData().
type RequiredProvidersData struct {
	AzAPIVersion   string
	AzureRMVersion string
}

const (
	requiredProvidersContent = `
terraform {
	required_version = ">= 1.0.0"
	required_providers {
		azurerm = {
			source  = "hashicorp/azurerm"
			version = "{{ .AzureRMVersion }}"
		}
		azapi = {
			source  = "azure/azapi"
			version = "{{ .AzAPIVersion }}"
		}
	}
}`
)

func generateRequiredProviders(data RequiredProvidersData, w io.Writer) error {
	tmpl := template.Must(template.New("terraformtf").Parse(requiredProvidersContent))
	return tmpl.Execute(w, data)
}

// GenerateRequiredProvidersFile generates a required providers file in the given path.
// The file path should be "terraform.tf".
func GenerateRequiredProvidersFile(data RequiredProvidersData, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return generateRequiredProviders(data, f)
}

// NewRequiredProvidersData generated a new version of the required providers data struct.
// It will use environment variables "AZAPI_VERSION" and "AZURERM_VERSION" to generate the data.
// If the environment variables are not set or the value is "latest", it will use the default values.
func NewRequiredProvidersData() RequiredProvidersData {
	var rpd RequiredProvidersData
	azapiver := ">= 0.3.0"
	azurermver := ">= 3.7.0"

	if val := os.Getenv("AZAPI_VERSION"); val != "" || val == "latest" {
		azapiver = "= " + val
	}
	if val := os.Getenv("AZURERM_VERSION"); val != "" || val == "latest" {
		azurermver = "= " + val
	}
	rpd.AzAPIVersion = azapiver
	rpd.AzureRMVersion = azurermver
	return rpd
}
