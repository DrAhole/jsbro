# JSBro

**JSBro** is a tool for scanning JavaScript endpoints for secrets by applying regex patterns defined in a YAML file.

## Features

- **Concurrent Processing:** Specify the number of concurrent HTTP requests.
- **Configurable Regex Patterns:** Provide a YAML file with regex patterns.
- **Custom Endpoint List:** Scan a user-provided file containing JS endpoint URLs.
- **Colored Output:** Easily distinguish findings in your terminal.
- **ASCII Banner:** Enjoy an attractive banner on startup.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/grumpzsux/jsbro.git
   cd jsbro
   go build -o jsbro main.go
```

2. **Install with Go:**
```
go install -v github.com/grumpzsux/jsbro@latest
```
