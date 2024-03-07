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
	Inferences []Inference `hcl:"infer,block"`
	Code       string      `hcl:"code,optional"`
}

type Inference struct {
	Assertion   string  `hcl:"assert"`
	Model       string  `hcl:"model"`
	Count       int     `hcl:"count"`
	Threshold   float64 `hcl:"threshold"`
	MaxTokens   int     `hcl:"max_tokens,optional"`
	Temperature float64 `hcl:"temperature,optional"`
	Tag_Name    string  `hcl:"tag_name,optional"`
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

	// Set the Tag_Name for each inference based on the parent tag's name
	for i := range config.Files {
		file := &config.Files[i]
		for j := range file.Tags {
			tag := &file.Tags[j]
			for k := range tag.Inferences {
				inference := &tag.Inferences[k]
				inference.Tag_Name = tag.Name
			}
		}
	}

	return &config, nil
}

func AttachCodeToTags(file *File) error {
	// Read the file content
	content, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err.Error())
	}

	// Convert the content to a string
	code := string(content)

	// Attach the code to each tag in the file
	for i := range file.Tags {
		file.Tags[i].Code = code
	}

	return nil
}
