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
			expectedError: `failed to parse Inferfile:
- ../testdata/missing_file_path.hcl:2,45-45: Missing required argument; The argument "code" is required, but no definition was found.
- ../testdata/missing_file_path.hcl:3,5-10: Unsupported block type; Blocks of type "infer" are not expected here.`,
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
