package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/inferret/infer/internal/api"
	"github.com/inferret/infer/internal/executor"
	"github.com/inferret/infer/internal/parser"
)

func customUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  validate: Validate the syntax of the Inferfile.")
	fmt.Println("  infer: Run inference tests (default).")
	fmt.Println("Options:")
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "openai-api-key" {
			fmt.Printf("  -%s string\n", f.Name)
			fmt.Printf("        %s (Default: environment variable $OPENAI_API_KEY)\n", f.Usage)
		} else {
			// Print the default help for other flags
			fmt.Printf("  -%s %s\n", f.Name, f.Usage+" (Default: "+f.DefValue+")")
		}
	})
}

func main() {
	// Infer: command line arguments
	flag.Usage = customUsage
	help := flag.Bool("help", false, "Show help information")
	apiKey := flag.String("openai-api-key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	apiUrl := flag.String("openai-api-url", "https://api.openai.com/v1", "OpenAI API URL")
	inferfile := flag.String("f", "Inferfile.hcl", "Path to the Inferfile")
	parallelThreads := flag.Int("parallel-threads", 1, "Number of parallel threads to run")
	verbose := flag.Bool("v", false, "Enable verbose output")
	// EndInfer: command line arguments
	flag.Parse()

	// Show help if -h or --help is provided
	if *help {
		customUsage()
		return
	}

	// Determine the command
	var command string
	if len(flag.Args()) > 0 {
		command = flag.Arg(0)
	} else {
		command = "infer" // Default command
	}

	switch command {
	case "validate":
		if err := validateInferfile(*inferfile); err != nil {
			fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Validation successful.")
	case "infer":
		if err := runInference(*apiKey, *apiUrl, *inferfile, *verbose, *parallelThreads); err != nil {
			fmt.Fprintf(os.Stderr, "Inference error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid command: %s\n", command)
		customUsage()
		os.Exit(1)
	}
}

func validateInferfile(filepath string) error {
	_, err := parser.ParseInferfile(filepath)
	return err
}

// Infer: execution
func runInference(apiKey, apiUrl, inferfilePath string, verbose bool, parallelThreads int) error {
	client := api.NewOpenAIWrapper(apiKey, apiUrl, verbose)
	config, err := parser.ParseInferfile(inferfilePath)
	if err != nil {
		return err
	}

	exec := executor.NewExecutor(client, verbose)
	errChan := make(chan error, len(config.Files))
	semaphore := make(chan struct{}, parallelThreads)

	var wg sync.WaitGroup

	for _, file := range config.Files {
		// Read the contents of the file
		code, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", file.Path, err)
		}

		for _, tag := range file.Tags {
			for _, inference := range tag.Inferences {
				wg.Add(1)
				go func(inference parser.Inference, code string) {
					defer wg.Done()
					semaphore <- struct{}{}
					if verbose {
						fmt.Printf("Inferring %s\n", inference.Assertion)
					}
					result, err := exec.Execute(inference, code)
					<-semaphore
					if err != nil {
						errChan <- err
					} else if !result {
						errChan <- fmt.Errorf("Inference failed: %s", inference.Assertion)
					} else if verbose {
						fmt.Printf("Inference successful: %s\n", inference.Assertion)
					}
				}(inference, string(code))
			}
		}
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
		return fmt.Errorf("Inference completed with errors")
	}

	return nil
}

// InferEnd: execution
