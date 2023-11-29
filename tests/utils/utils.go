package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
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
		return logger.TestingT
	}
	return logger.Discard
}

// PreCheckDeployTests ensures the correct environment variables
// are set for the deployment tests to run.
func PreCheckDeployTests(t *testing.T) {
	// Skip if we haven't set the `TERRATEST_DEPLOY` variable.
	if value := os.Getenv("TERRATEST_DEPLOY"); value == "" {
		t.Skip("`TERRATEST_DEPLOY` must be set to `true` for deployment tests! - Skipping...")
	}
	// These variables cause a failure if not set.
	variables := []string{
		"AZURE_BILLING_SCOPE",
		"AZURE_TENANT_ID",
		"AZURE_SUBSCRIPTION_ID",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Logf("`%s` must be set for deployment tests!", variable)
			t.FailNow()
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

// GetTestDir returns the directory of the test file.
func GetTestDir(t *testing.T) string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
