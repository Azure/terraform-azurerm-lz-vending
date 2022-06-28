package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"gopkg.in/matryer/try.v1"
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
		return logger.Default
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
		"AZURE_SUBSCRIPTION_ID",
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

// GetDefaultTerraformOptions returns the default Terraform options for the
// given directory.
func GetDefaultTerraformOptions(t *testing.T, dir string) *terraform.Options {
	if dir == "" {
		dir = "./"
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	pf := dir + "tfplan"
	return &terraform.Options{
		Logger:       GetLogger(),
		NoColor:      true,
		PlanFilePath: pf,
		TerraformDir: dir,
		Vars:         make(map[string]interface{}),
	}
}

// GetTestDir returns the directory of the test file.
func GetTestDir(t *testing.T) string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

// TerraformDestroyWithRetry is a helper function that wraps a terraform destroy in a try.Do
// designed to be used as a defer function.
func TerraformDestroyWithRetry(t *testing.T, to *terraform.Options, dur time.Duration, max int) {
	if try.MaxRetries < max {
		try.MaxRetries = max
	}
	err := try.Do(func(attempt int) (bool, error) {
		_, err := terraform.DestroyE(t, to)
		if err != nil {
			t.Logf("terraform destroy failed, attempt %d/%d", attempt, max)
			time.Sleep(dur)
		}
		return attempt < max, err
	})
	if err != nil {
		t.Logf("terraform destroy error: %v", err)
	}
}

// removeTestDir removes the supplied test directory
func removeTestDir(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Logf("Error removing test directory %s: %v", dir, err)
	}
}

// CopyTerraformFolderToTempAndCleanUp sets up a temporary copy of the supplied module folder
// It will return the temporary directory path and a function that can be used to clean up the temporary directory
//
// The testdir input is the relative path to the test directory, it can be blank if testing the module directly with variables
// or it can be a relative path to the module directory if testing the module using a subdirectory.
//
// Note: This function will only work if the test directory is in a child subdirectory of the test directory.
// e.g. you cannot use parent paths of the moduleDir.
//
// The depth input is used to determine how many directories to go up to make sure we
// fully clean up.
//
// The function will return the temporary directory to use with the terraform options struct, as well as
//  a function that can be used with defer to clean up afterwards.
func CopyTerraformFolderToTempAndCleanUp(t *testing.T, moduleDir string, testDir string) (string, func(), error) {
	tmp := test_structure.CopyTerraformFolderToTemp(t, moduleDir, testDir)
	// We normalise, then work out the depth of the test directory relative
	// to the test so we know how many/ directories to go up to get to the root.
	// We can then delete the right directory when cleaning up.
	absTestPath := filepath.Join(moduleDir, testDir)
	relPath, err := filepath.Rel(moduleDir, absTestPath)
	if err != nil {
		return "", nil, fmt.Errorf("could not get relative path to test directory: %v", err)
	}
	list := strings.Split(relPath, string(os.PathSeparator))
	depth := len(list)
	if len(list) > 1 {
		depth++
	}
	dir := tmp
	f := func() {
		for i := 0; i < depth; i++ {
			dir = filepath.Dir(dir)
		}
		removeTestDir(t, dir)
	}
	return tmp, f, nil
}

// CreateTerraformProvidersFile creates an azurerm terraform providers file in the supplied directory.
func CreateTerraformProvidersFile(dir string) error {
	dir = filepath.Clean(dir)
	f, err := os.Create(filepath.Join(dir, "providers.tf"))
	if err != nil {
		return fmt.Errorf("error creating providers.tf: %v", err)
	}
	providerstf := `
provider "azurerm" {
  features {}
}`
	_, err = f.WriteString(providerstf)
	if err != nil {
		return fmt.Errorf("error writing providers.tf: %v", err)
	}
	return nil
}
