# go-envfile

`go-envfile` is a simple and easy-to-use Go library for loading environment variables from `.env` files. It supports automatically loading configuration files for different environments based on the `GO_ENV` environment variable and provides a sensible default loading order.

## Features

* **Environment Aware:** Automatically loads corresponding `.env` files based on the `GO_ENV` environment variable.
* **Multiple Environment Support by Default:** Supports `development`, `production`, and `test` environments by default.
* **File Loading Priority:** Defines a clear file loading priority for each environment, making environment overriding easy.
* **Local Overrides:** Supports the `.local` suffix for local development environment overrides.
* **Simple to Use:** Provides a single `Load()` function to load environment variables.
* **Error Handling:** Provides detailed log output for easy tracking of errors during file loading and environment variable setting.
* **Comment and Empty Line Ignoring:** Automatically ignores comments (starting with `#`) and empty lines in `.env` files.
* **Basic Format Validation:** Checks for the basic `key=value` format on each line.

## Installation

You can install `go-envfile` using the `go get` command:

```bash
go get https://github.com/lucap9056/go-envfile/envfile
```

## Usage

Import the `envfile` package into your Go code and call the `envfile.Load()` function early in your program to load environment variables.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"https://github.com/lucap9056/go-envfile/envfile"
)

func main() {
	// Load .env files
	envfile.Load()

	// Now you can retrieve environment variables using os.Getenv()
	apiKey := os.Getenv("API_KEY")
	databaseURL := os.Getenv("DATABASE_URL")

	fmt.Println("API_KEY:", apiKey)
	fmt.Println("DATABASE_URL:", databaseURL)

	// ... your application logic ...
}
```

### Example `.env` File

```
# This is a comment
$user=user
$PSWD=password
API_KEY=your_api_key_here
DATABASE_URL=postgres://{$user}:{$PSWD}@host:port/database
PORT=8080

# Empty lines are ignored

DEBUG=true
```

### Environment-Specific `.env` Files

`go-envfile` determines which `.env` files to load based on the value of the `GO_ENV` environment variable. If `GO_ENV` is not set or is set to an unrecognized value, it defaults to loading configuration files for the `development` environment.

The following is the loading priority for different environments:

**Development:**

1.  `.env.development.local`
2.  `.env.dev.local`
3.  `.env.development`
4.  `.env.dev`
5.  `.env.local`
6.  `.env`

**Production:**

1.  `.env.production.local`
2.  `.env.prod.local`
3.  `.env.production`
4.  `.env.prod`
5.  `.env.local`
6.  `.env`

**Test:**

1.  `.env.test.local`
2.  `.env.test`
3.  `.env.testing`
4.  `.env.local`
5.  `.env`

Once `go-envfile` finds an existing file and loads it successfully, it stops searching and returns. This means that files listed earlier have a higher priority and can override settings in later files.

### Local Overrides

Files with the `.local` suffix (e.g., `.env.development.local` or `.env.local`) are typically used for local development environments. Settings in these files will override settings in the corresponding files without the `.local` suffix. This allows developers to have different configurations on their local machines without modifying the main `.env` files.

### Log Output

`go-envfile` uses the `log` package to output information and errors during the loading process:

* **Warning:** Output when the `GO_ENV` environment variable is not recognized or when no `.env` file is successfully loaded.
* **Error:** Output when there is a failure to open a file, the file format is incorrect, or setting an environment variable fails.
* **Info:** Output when a `.env` file is successfully loaded, indicating the path of the loaded file.

## Important Notes

* Ensure that your `.env` files are located in the current working directory of your application.
* Do not commit sensitive information (e.g., passwords, API keys) directly into your version control system. It is recommended to add `.env` files to your `.gitignore`.
* In production environments, it is generally recommended to manage environment variables through more secure methods, such as system environment variables or dedicated configuration management tools. `.env` files are more suitable for development and testing environments.

## References

* **godotenv:** [https://github.com/joho/godotenv](https://github.com/joho/godotenv) - A widely used Go library for loading `.env` files that provided some inspiration for the design of `go-envfile`.
