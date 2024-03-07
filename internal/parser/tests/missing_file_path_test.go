package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/inferret/infer/internal/parser"
)

func TestParseInferfileMissingFilePath(t *testing.T) {
	testCases := []struct {
		name          string
		inferfilePath string
		expectedErr   string
	}{
		{
			name:          "missing_file_path",
			inferfilePath: filepath.Join("testdata", "missing_file_path.hcl"),
			expectedErr:   "file not found: ./this/file/does/not/exist (Inferfile line: 1)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.ParseInferfile(tc.inferfilePath)
			if err == nil || err.Error() != tc.expectedErr {
				t.Errorf("Expected error: %s, got: %v", tc.expectedErr, err)
			}
		})
	}
}
