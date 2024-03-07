package parser_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/inferret/infer/internal/parser"
)

func TestParseInferfileMissingFilePath(t *testing.T) {
	testCases := []struct {
		name          string
		inferfilePath string
		expectedError string
	}{
		{
			name:          "missing_file_path",
			inferfilePath: filepath.Join("../testdata", "missing_file_path.hcl"),
			expectedError: `file '/this/file/does/not/exist' specified in Inferfile does not exist`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.ParseInferfile(tc.inferfilePath)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if strings.TrimSpace(err.Error()) != strings.TrimSpace(tc.expectedError) {
				t.Errorf("Expected error:\n%s\ngot:\n%s", tc.expectedError, err.Error())
			}
		})
	}
}
