package envfile

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Load reads environment variables from a list of potential .env files.
// It prioritizes files based on the current environment specified by the
// "GO_ENV" environment variable. If "GO_ENV" is not set or invalid, it
// defaults to loading from development-related .env files.
// It searches for the following files in the current directory, in order
// of precedence for each environment:
//
// Development:
//
// .env.development.local, .env.dev.local, .env.development, .env.dev, .env.local, .env
//
// Production:
//
// .env.production.local, .env.prod.local, .env.production, .env.prod, .env.local, .env
//
// Test:
//
// .env.test.local, .env.test, .env.testing, .env.local, .env
//
// If a file is found and successfully loaded, the function returns. If
// errors occur during file reading or environment variable setting, they
// are logged. A warning is logged if no .env file is successfully loaded.
func Load() {
	envFileMap := map[string][]string{
		"development": {
			".env.development.local",
			".env.dev.local",
			".env.development",
			".env.dev",
			".env.local",
			".env",
		},
		"production": {
			".env.production.local",
			".env.prod.local",
			".env.production",
			".env.prod",
			".env.local",
			".env",
		},
		"test": {
			".env.test.local",
			".env.test",
			".env.testing",
			".env.local",
			".env",
		},
	}

	env := os.Getenv("GO_ENV")
	envNames, exists := envFileMap[env]
	if !exists {
		log.Printf("Warning: Environment '%s' is not recognized. Defaulting to 'development' environment files.", env)
		envNames = envFileMap["development"]
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error: Could not get the current working directory: %v", err)
		return
	}

	files, err := os.ReadDir(cwd)
	if err != nil {
		log.Printf("Error: Could not read the current directory: %v", err)
		return
	}

	fileMap := make(map[string]struct{})
	for _, file := range files {
		if !file.IsDir() {
			fileMap[file.Name()] = struct{}{}
		}
	}

	loaded := false

	for _, name := range envNames {
		if _, exists := fileMap[name]; exists {

			filePath := filepath.Join(cwd, name)

			if err := loadFile(filePath); err != nil {
				log.Printf("Error: Failed to load environment variables from '%s': %v", filePath, err)
			} else {
				log.Printf("Successfully loaded environment variables from '%s'", filePath)
				loaded = true
				return
			}
		}
	}

	if !loaded {
		log.Println("Warning: No .env file was successfully loaded. Ensure at least one of the expected .env files exists in the current directory.")
	}
}

func loadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error: unable to open file '%s': %v", filePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error: Failed to close file '%s': %v", filePath, err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Warning: Invalid line format in '%s': '%s'. Expected 'key=value'.", filePath, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			log.Printf("Warning: Empty key found in '%s': '%s'. Skipping.", filePath, line)
			continue
		}

		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("error: unable to set environment variable '%s': %v", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error: failed to read file '%s': %v", filePath, err)
	}

	return nil
}
