package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// ReadEnv reads an env file and returns a map of the key-value pairs
func ReadEnv(filePath string) (map[string]string, error) {
	envMap := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}

// MergeEnvFiles merges multiple env files, later files overriding earlier ones
func MergeEnvFiles(outputFile string, inputFiles []string) error {
	finalEnv := make(map[string]string)

	for _, inputFile := range inputFiles {
		env, err := ReadEnv(inputFile)
		if err != nil {
			return err
		}

		// Merge env maps, later files override earlier ones
		for key, value := range env {
			finalEnv[key] = value
		}
	}

	// Write the merged env to the output file
	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	for key, value := range finalEnv {
		_, err := fmt.Fprintf(output, "%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Define a flag for the output file
	outputFile := flag.String("o", "merged.env", "Specify the output file name")
	flag.Parse()

	// The remaining arguments are the input files
	inputFiles := flag.Args()

	// Check if there are enough input files
	if len(inputFiles) == 0 {
		fmt.Println("Usage: merge_env -o <output-file> <input-files...>")
		return
	}

	err := MergeEnvFiles(*outputFile, inputFiles)
	if err != nil {
		fmt.Printf("Error merging env files: %v\n", err)
	}
}
