package parser_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/inferret/infer/internal/parser"
)

func TestParseInferfileValid(t *testing.T) {
	testCases := []struct {
		name           string
		inferfilePath  string
		expectedConfig *parser.InferConfiguration
	}{
		{
			name:          "valid",
			inferfilePath: filepath.Join("../testdata", "valid.hcl"),
			expectedConfig: &parser.InferConfiguration{
				Files: []parser.File{
					{
						Path: filepath.Join("../testdata", "valid.go"),
						Tags: []parser.Tag{
							{
								Name: "OpenAI client",
								Inferences: []parser.Inference{
									{
										Assertion:   "Does not contain any hard-coded credentials.",
										Model:       "gpt-3.5-turbo",
										Count:       5,
										Threshold:   0.8,
										MaxTokens:   0,
										Temperature: 0,
										Tag_Name:    "OpenAI client",
									},
								},
								Code: `func openaiClient() {
	// This is a stub function
}
`,
							},
						},
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

			// Attach code to tags
			for i := range config.Files {
				err := parser.AttachCodeToTags(&config.Files[i])
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			if !reflect.DeepEqual(config, tc.expectedConfig) {
				t.Errorf("Expected config: %+v, got: %+v", tc.expectedConfig, config)
			}
		})
	}
}
