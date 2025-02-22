// main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Pattern represents a regex pattern entry.
type Pattern struct {
    Name       string `yaml:"name"`
    Regex      string `yaml:"regex"`
    Confidence string `yaml:"confidence,omitempty"`
}

// YamlPattern wraps the Pattern field in the YAML file.
type YamlPattern struct {
    Pattern Pattern `yaml:"pattern"`
}

// Config is the top-level structure of your YAML file.
type Config struct {
    Patterns []YamlPattern `yaml:"patterns"`
}

// Global CLI flags
var (
	listFile    string
	configFile  string
	concurrency int
	timeout     int
)

// asciiBanner is the ASCII art banner displayed on startup.
const asciiBanner = `

     _________         
 __ / / __/ _ )_______ 
/ // /\ \/ _  / __/ _ \
\___/___/____/_/  \___/ [v1.0]

 GRuMPzSux | www.grumpz.net

[+] Scan JS endpoints for secrets using regex patterns
                                  
==================================================================        
`

func init() {
	rootCmd.Flags().StringVarP(&listFile, "list", "l", "", "File containing JS endpoint URLs (one per line) (required)")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "YAML config file containing regex patterns (required)")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "n", 5, "Number of concurrent requests")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "HTTP request timeout in seconds")
	_ = rootCmd.MarkFlagRequired("list")
	_ = rootCmd.MarkFlagRequired("config")
}

var rootCmd = &cobra.Command{
    Use:   "jsbro",
    Short: "JSBro - A tool for scanning JS endpoints for secrets using regex patterns",
    PreRun: func(cmd *cobra.Command, args []string) {
        fmt.Println(asciiBanner)
    },
    Run: func(cmd *cobra.Command, args []string) {
        runJSBro()
    },
}


func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runJSBro is the main execution function.
func runJSBro() {
	// Load the regex patterns from the YAML configuration.
	patterns, err := loadPatterns(configFile)
	if err != nil {
		log.Fatalf("Failed to load regex patterns: %v", err)
	}
	color.Green("[*] Loaded %d pattern(s) from config", len(patterns))

	// Load JS endpoints from the provided file.
	endpoints, err := loadEndpoints(listFile)
	if err != nil {
		log.Fatalf("Failed to load JS endpoints: %v", err)
	}
	color.Green("[*] Loaded %d JS endpoint(s)", len(endpoints))

	// Prepare HTTP client with timeout.
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Semaphore channel to control concurrency.
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, url := range endpoints {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			processEndpoint(u, patterns, client)
			<-sem // Release semaphore
		}(url)
	}
	wg.Wait()
	color.Green("[*] Scanning complete!")
}

// loadPatterns reads and parses the YAML config file into a slice of Pattern.
func loadPatterns(filename string) ([]Pattern, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    var patterns []Pattern
    for _, yp := range config.Patterns {
        patterns = append(patterns, yp.Pattern)
    }
    return patterns, nil
}

func loadEndpoints(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var endpoints []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			endpoints = append(endpoints, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return endpoints, nil
}

// processEndpoint fetches the JS file at the given URL and applies all regex patterns.
func processEndpoint(url string, patterns []Pattern, client *http.Client) {
	color.Cyan("[*] Fetching: %s", url)
	resp, err := client.Get(url)
	//if err != nil {
	//	color.Red("[-] Error fetching %s: %v", url, err)
	//	return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red("[-] Non-OK HTTP status for %s: %s", url, resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red("[-] Error reading body for %s: %v", url, err)
		return
	}

	// For each pattern, search the JS content.
	for _, pat := range patterns {
		regex, err := regexp.Compile(pat.Regex)
		if err != nil {
			color.Red("[-] Invalid regex for pattern '%s': %v", pat.Name, err)
			continue
		}

		matches := regex.FindAllString(string(body), -1)
		if len(matches) > 0 {
			// For each match, print a result line.
			for _, match := range matches {
				// Format: [+] [Name Found] - [Found Data Goes Here] - [Full URL]
				// We'll use different colors for clarity.
				color.Green("[+] %s", pat.Name)
				color.Yellow("    Data: %s", match)
				color.Blue("    URL : %s", url)
			}
		}
	}
}
