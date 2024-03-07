package parser_test

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestMain is the entry point for all tests in this package.
// It sets up any necessary global state and runs the tests.
func TestMain(m *testing.M) {
	// Perform any global setup here, if needed.
	// For example, you can create temporary directories or files.

	// Run the tests.
	exitCode := m.Run()

	// Perform any global teardown here, if needed.
	// For example, you can clean up temporary directories or files.

	// Exit with the appropriate code.
	os.Exit(exitCode)
}

// createStubFile creates a temporary stub Go file for testing purposes.
func createStubFile(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "stub*.go")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
