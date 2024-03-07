package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type InferConfiguration struct {
	Files []File `hcl:"file,block"`
}

type File struct {
	Path string `hcl:"path,label"`
	Tags []Tag  `hcl:"tag,block"`
}

type Tag struct {
	Name       string      `hcl:"name,label"`
	Inferences []Inference `hcl:"inference,block"`
	Code       string      `hcl:"code"`
}

type Inference struct {
	Assertion   string  `hcl:"assertion"`
	Model       string  `hcl:"model"`
	Count       int     `hcl:"count"`
	Threshold   float64 `hcl:"threshold"`
	MaxTokens   int     `hcl:"max_tokens"`
	Temperature float64 `hcl:"temperature"`
	Tag_Name    string  `hcl:"tag_name"`
}

func ParseInferfile(path string) (*InferConfiguration, error) {
	var config InferConfiguration

	// Read the Inferfile
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Inferfile: %s", err.Error())
	}

	// Parse the Inferfile
	err = hclsimple.Decode(path, src, nil, &config)
	if err != nil {
		// Extract the detailed error message from the hcl.Diagnostics
		var errMsg strings.Builder
		errMsg.WriteString("failed to parse Inferfile:\n")
		for _, diag := range err.(hcl.Diagnostics) {
			errMsg.WriteString(fmt.Sprintf("- %s\n", diag.Error()))
		}
		return nil, fmt.Errorf(errMsg.String())
	}

	// Get the directory of the Inferfile
	inferfileDir := filepath.Dir(path)

	// Update file paths to be relative to the Inferfile's directory
	for i := range config.Files {
		file := &config.Files[i]
		if !filepath.IsAbs(file.Path) {
			file.Path = filepath.Join(inferfileDir, file.Path)
		}
		// Check if the file exists
		if _, err := os.Stat(file.Path); os.IsNotExist(err) {
			return nil, fmt.Errorf("file '%s' specified in Inferfile does not exist", file.Path)
		}
	}

	return &config, nil
}

func AttachCodeToTags(file *File) error {
	// Read the entire source file content
	source, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return err
	}

	// Convert the file content into a slice of lines
	lines := strings.Split(string(source), "\n")

	// Iterate over each tag in the file
	for i := range file.Tags {
		tag := &file.Tags[i] // Get a reference to the tag to modify it directly
		var tagBuilder strings.Builder
		var inTagBlock bool

		// Iterate over each line in the source file
		for _, line := range lines {
			if strings.Contains(line, "Infer: "+tag.Name) {
				inTagBlock = true // Start of tag block
				continue
			}
			if strings.Contains(line, "EndInfer: "+tag.Name) {
				inTagBlock = false // End of tag block
				break
			}
			if inTagBlock {
				tagBuilder.WriteString(line + "\n") // Collect the lines within the tag block
			}
		}

		// Update the tag's Code field with the collected code block
		tag.Code = tagBuilder.String()
	}

	return nil
}
