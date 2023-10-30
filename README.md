# Gosha: Go SHA256 Hash Generator for Packages

[![Go Report Card](https://goreportcard.com/badge/github.com/draganm/gosha)](https://goreportcard.com/report/github.com/draganm/gosha)
[![GoDoc](https://pkg.go.dev/badge/github.com/draganm/gosha)](https://pkg.go.dev/github.com/draganm/gosha)
![License](https://img.shields.io/github/license/draganm/gosha)

`Gosha` is a versatile Go package and accompanying CLI tool designed to generate SHA256 hashes for Go packages and their dependencies. This becomes invaluable for integrity checks in CI/CD pipelines, automated workflows, or even managing monorepos.

## 🌟 Key Features

- 📦 **Package Hashing**: Generate SHA256 hashes for any Go package.
- 🛠️ **CLI & Library Support**: Both command-line and programmatic interfaces are available.
- ⚙️ **Fine-grained Control**: Optionally include standard library and test files in the hash generation.
- 🚀 **Use-cases**: 
  - Efficiently manage monorepo builds by rebuilding only when a service source code has changed.
  - Use the generated SHA as a tag for Docker images, ensuring Kubernetes Deployments are updated only when necessary.

## 📥 Installation

Install the package and CLI tool using `go get`:

```bash
go get -u github.com/draganm/gosha
```

## 📘 Usage

### CLI Interface

Use the following command syntax:

```bash
gosha [OPTIONS] [PACKAGE_PATH]
```

**Options**:

- `--include-stdlib`: Include Go's standard libraries in the hash generation.  
  - Environment variable: `INCLUDE_STDLIB`
- `--include-testfiles`: Include test files in the hash generation.  
  - Environment variable: `INCLUDE_TESTFILES`

#### Examples:

To generate a SHA256 hash for the package in the current directory:

```bash
gosha
```

To include standard libraries for a specific package:

```bash
gosha --include-stdlib <path to your main package>
```

### Programmatic Interface

To use Gosha programmatically, import the `gosha` package and call the `CalculatePackageSHA()` function.

```go
import "github.com/draganm/gosha/gosha"

finalSHA, err := gosha.CalculatePackageSHA("<path to your main package>", false, false)
if err != nil {
    fmt.Println("Error:", err)
    return
}
// Use the finalSHA as needed...
```

## 👥 Contributing

Contributions are welcome! Feel free to submit issues for bug reports, feature requests, or even pull requests.

## 📜 License

This project is licensed under the MIT License.
