package parser_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/inferret/infer/internal/parser"
)

func TestParseInferfileNoTags(t *testing.T) {
	testCases := []struct {
		name           string
		inferfilePath  string
		expectedConfig *parser.InferConfiguration
	}{
		{
			name:          "no_tags",
			inferfilePath: filepath.Join("../testdata", "no_tags.hcl"),
			expectedConfig: &parser.InferConfiguration{
				Files: []parser.File{
					{
						Path: "/tmp/mycode.go",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := parser.ParseInferfile(tc.inferfilePath)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(config, tc.expectedConfig) {
				t.Errorf("Expected config: %+v, got: %+v", tc.expectedConfig, config)
			}
		})
	}
}
