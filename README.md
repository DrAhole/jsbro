# JSBro

![JSBro Logo](./assets/jsbro-logo.png) <!-- Replace with your logo or banner image -->

JSBro is a powerful and fast tool written in Go for scanning JavaScript endpoints for secrets and sensitive information. By leveraging user-defined regex patterns (via YAML configuration), JSBro helps you quickly identify potential security leaks across a list of JavaScript URLs.

## Features

- **Concurrent Processing:** Scan multiple JS endpoints at once with configurable concurrency.
- **Customizable Regex Patterns:** Easily define and update regex patterns using a YAML configuration file.
- **Colorful, User-Friendly Output:** Results are clearly presented in your terminal with color coding for easy identification.
- **ASCII Art Banner:** Enjoy a stylish ASCII banner on startup.

## Installation

Make sure you have [Go](https://golang.org/dl/) installed (version 1.16+ recommended).

You can install JSBro directly using the `go install` command:

```bash
go install -v github.com/grumpzsux/jsbro@latest
```
This will compile JSBro and install the binary into your `$GOPATH/bin`.

Alternatively, clone the repository and build it manually:
```bash
git clone https://github.com/grumpzsux/jsbro.git
cd jsbro
go build -o jsbro main.go
```
## Usage
JSBro requires two inputs:
- A list file (`--list` or `-l`) that contains one JavaScript endpoint URL per line.
- A YAML configuration file (`--config` or `-c`) that defines the regex patterns to search for, check the `/patterns/` directory.

**Example command:**
```bash
./jsbro --list /path/to/endpoints.txt --config /path/to/patterns.yaml --concurrency 10 --timeout 10
```
## Example YAML Configuration

Below is an example of a YAML configuration file:
```yaml
patterns:
  - pattern:
      name: AWS Access Key
      regex: "(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}"
      confidence: high
  - pattern:
      name: AWS Secret Key
      regex: "(?i)aws(.{0,20})?(?-i)['\\\"][0-9a-zA-Z\\/+]{40}['\\\"]"
      confidence: high
```

## Contributing
Contributions are welcome! Please fork the repository and submit your pull requests. If you find any issues or have suggestions, feel free to open an issue on GitHub.

## Contact
For questions or support, send me a Direct Message on X @GRuMPzSux
