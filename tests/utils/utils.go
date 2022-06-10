package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const (
	planFilePath = "../tfplan"
	terraformDir = "../"
)

// SanitiseErrorMessage replaces the newline characters in an error.Error() output with a single space to allow us to check for the entire error message.
// We need to do this because Terraform adds newline characters depending on the width of the console window.
// TODO: Test on Windows if we get \r\n instead of just \n.
func SanitiseErrorMessage(err error) string {
	return strings.Replace(err.Error(), "\n", " ", -1)
}

// GetLogger returns a logger that can be used for testing.
// The default logger will discard the Terraform output.
// Set TERRATEST_LOGGER to a non empty value to enable verbose logging.
func GetLogger() *logger.Logger {
	if os.Getenv("TERRATEST_LOG") != "" {
		return logger.Terratest
	}
	return logger.Discard
}

// PreCheckDeployTests ensures the correct environment variables
// are set for the deployment tests to run.
func PreCheckDeployTests(t *testing.T) {
	variables := []string{
		"TERRATEST_DEPLOY",
		"AZURE_BILLING_SCOPE",
		"AZURE_TENANT_ID",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for deployment tests!", variable)
		}
	}
}

// RandomHex generates a random hex string of the given byte length.
// Uses crypto/rand for generating the random bytes not math/rand
// as we kept getting the same results from the math/rand generator.
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetDefaultTerraformOptions returns a TerraformOptions struct with the correct values
func GetDefaultTerraformOptions(vars map[string]interface{}) *terraform.Options {
	return &terraform.Options{
		TerraformDir: terraformDir,
		NoColor:      true,
		Logger:       GetLogger(),
		PlanFilePath: planFilePath,
		Vars:         vars,
	}
}
